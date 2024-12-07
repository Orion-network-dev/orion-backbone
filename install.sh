#!/usr/bin/env bash

# Script: Orion V3 Installation Script
# Description: Downloads and installs the latest version of OrionV3 for the current architecture
# Usage: sudo ./install_orion.sh
# Return codes:
#   0: Success
#   1: General error
#   2: Root privileges required
#   3: Network error
#   4: Package installation error

# Enable strict error handling
set -euo pipefail

# Constants
readonly GITHUB_API="https://api.github.com/repos/MatthieuCoder/OrionV3/releases/latest"
readonly CURL_TIMEOUT=30
readonly CURL_RETRIES=3
readonly SERVICE_NAME="oriond"

# Logging functions
log_info() {
    echo "[INFO] $1"
}

log_error() {
    echo "[ERROR] $1" >&2
}

log_warning() {
    echo "[WARNING] $1" >&2
}

# Cleanup function
cleanup() {
    local temp_file="${TEMP_DEB:-}"
    if [[ -f "${temp_file}" ]]; then
        log_info "Cleaning up temporary files..."
        rm -f "${temp_file}"
    fi
}

# Error handler
error_handler() {
    local line_no=$1
    local error_code=$2
    log_error "Error occurred in script at line: ${line_no}, with error code: ${error_code}"
    cleanup
    exit 1
}

# Set up error handling
trap 'error_handler ${LINENO} $?' ERR
trap cleanup EXIT

# Check for root privileges
if [[ $EUID -ne 0 ]]; then
    log_error "This script must be run as root"
    exit 2
}

# Verify and install dependencies
install_dependency() {
    local package=$1
    if ! command -v "${package}" &> /dev/null; then
        log_info "The ${package} command couldn't be found in your current ${PATH}... Installing it."
        if ! apt-get install -yq "${package}"; then
            log_error "Failed to install ${package}"
            exit 4
        fi
    fi
}

# Update package lists
update_apt() {
    log_info "Running apt-get update for accurate package information..."
    if ! apt-get update -q; then
        log_error "Failed to update apt package lists"
        exit 4
    fi
}

# Download file with retries
download_with_retry() {
    local url=$1
    local output=$2
    local retry_count=0

    while [[ ${retry_count} -lt ${CURL_RETRIES} ]]; do
        if curl --fail \
            --silent \
            --location \
            --connect-timeout "${CURL_TIMEOUT}" \
            --retry 3 \
            --output "${output}" \
            "${url}"; then
            return 0
        fi
        retry_count=$((retry_count + 1))
        log_warning "Download failed, attempt ${retry_count}/${CURL_RETRIES}"
        sleep 2
    done
    return 1
}

# Verify service exists
verify_service() {
    if ! systemctl list-unit-files "${SERVICE_NAME}.service" &> /dev/null; then
        log_error "Service ${SERVICE_NAME} does not exist"
        exit 1
    fi
}

# Main execution
main() {
    # Create temporary file
    TEMP_DEB=$(mktemp)

    # Update package lists and install dependencies
    update_apt
    install_dependency "jq"
    install_dependency "curl"

    # Get latest release information
    log_info "Fetching latest release information..."
    local json
    json=$(curl -sf --connect-timeout "${CURL_TIMEOUT}" "${GITHUB_API}") || {
        log_error "Failed to fetch release information"
        exit 3
    }

    # Parse version and URL
    local version
    version=$(echo "${json}" | jq -r '.name') || {
        log_error "Failed to parse version information"
        exit 1
    }

    local arch
    arch=$(dpkg --print-architecture)
    local name_predicate="contains(\"${arch}.deb\")"
    local url
    url=$(echo "${json}" | jq -r ".assets[] | select(.name | ${name_predicate}) | .browser_download_url") || {
        log_error "Failed to find package for architecture: ${arch}"
        exit 1
    }

    if [[ -z "${url}" ]]; then
        log_error "No package found for architecture: ${arch}"
        exit 1
    }

    # Download package
    log_info "Downloading version ${version} for ${arch}..."
    if ! download_with_retry "${url}" "${TEMP_DEB}"; then
        log_error "Failed to download package"
        exit 3
    fi

    # Verify package
    if ! dpkg-deb -I "${TEMP_DEB}" &> /dev/null; then
        log_error "Invalid package file"
        exit 4
    }

    # Install package
    log_info "Installing version ${version}..."
    if ! apt-get install -q --allow-downgrades -y "${TEMP_DEB}"; then
        log_error "Failed to install package"
        exit 4
    fi

    # Verify and restart service
    verify_service
    log_info "Restarting services..."
    if ! systemctl restart "${SERVICE_NAME}"; then
        log_error "Failed to restart ${SERVICE_NAME}"
        exit 1
    fi

    # Wait for service to start
    sleep 2
    if ! systemctl is-active --quiet "${SERVICE_NAME}"; then
        log_error "Service failed to start"
        exit 1
    }

    log_info "Installation completed successfully"
}

# Execute main function
main

exit 0
