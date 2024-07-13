rm *.{crt,key}
step-cli ca root ca.crt
step-cli ca certificate "puffer.fish" server.crt server.key
step-cli ca certificate "0.member.orionet.re" client0.crt client0.key
step-cli ca certificate "1.member.orionet.re" client1.crt client1.key
