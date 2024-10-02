# Golang user servise


## Technologies:

+ Golang (net/http)
+ Router (chi)
+ PostgreSQL
+ Ldap (For Active Directory connection)

Environs:

1. **storage_path** - Database link

2. **secret_key** -  Generated secret_key for JWT authorization

3. **http_server** - Configuration of your http server

4. **ldap** - Configuration of ldap connection

    + *base_dn* - Base search group "DC=example, DC=com"
    + *bind_dn* - Admin username of Active Directory
    + *bind_password* - Admin password of Active Directory
    + *ldap_host* - Host of ldap server without ldap:// or ldaps://
    + *ldap_port* - 389 for ldap:// and 636 for ldaps://
    + *groupDN* - The group to be searched for (optional)

  ## Run flags
  **CONFIG_PATH** =. /internal/config/development.yaml
