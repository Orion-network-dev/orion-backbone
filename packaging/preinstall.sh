dpkg-divert --package orion-backbone --add --rename \
    --divert /etc/frr/daemons.original /etc/frr/daemons
