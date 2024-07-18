platform-config-root-dir: 'platform-config'
environments:
  - name: 'snowflake-hacks'
    snowflake:
      account: '<snowflake account>'
      user: '<snowflake user>'
      private-key-path: 'secrets/rsa_key.p8'
      role: '<snowflake role>'
      region: '<snwoflake region>'
      warehouse: '<snowflake warehouse'
      default-database: '<snowflake database>'
