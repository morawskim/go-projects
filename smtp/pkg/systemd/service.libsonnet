local envFile = {
  debian: "/etc/default/gosmtp",
  rpm: "/etc/sysconfig/gosmtp",
};

local service(packager) = |||
 [Unit]
 Description=A SMTP server that receives emails and forwards them as notifications
 Requires=network-online.target
 After=network-online.target

 [Service]
 Type=simple
 User=daemon
 Restart=always
 RestartSec=30
 EnvironmentFile=%s
 ExecStart=/usr/bin/gosmtp $ARGS

 [Install]
 WantedBy=multi-user.target
||| % envFile[packager];

{
  "gosmtp-deb.service": service("debian"),
  "gosmtp-rpm.service": service("rpm"),
}
