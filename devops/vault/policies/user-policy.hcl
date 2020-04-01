path "transit/*"
{
    capabilities = ["create", "update", "read"]
}

path "secret/*"
{
    capabilities = ["create", "update", "read"]
}

path "secret/gitea/*"
{
    capabilities = ["read"]
}
