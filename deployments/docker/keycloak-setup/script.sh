#!/bin/sh

${HOME}/bin/kcadm.sh config credentials  --server ${KEYCLOAK_URL} --realm master --user ${KEYCLOAK_ADMIN} --password ${KEYCLOAK_ADMIN_PASSWORD}
${HOME}/bin/kcadm.sh create realms -s realm=enricherrealm -s enabled=true
CID=$(${HOME}/bin/kcadm.sh create clients -r enricherrealm -f /keycloak-setup/userClient.json --id)
${HOME}/bin/kcadm.sh create users -r enricherrealm -f /keycloak-setup/user.json
${HOME}/bin/kcadm.sh set-password -r enricherrealm --username=testuser --new-password=userpassword
${HOME}/bin/kcadm.sh create clients/${CID}/roles -r enricherrealm -s name=userrole -s description="Роль для пользователя"
${HOME}/bin/kcadm.sh add-roles -r enricherrealm --uusername testuser --cclientid userClient --rolename userrole