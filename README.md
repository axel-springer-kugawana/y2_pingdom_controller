# y2_pingdom_controller

[![CI](https://github.com/axel-springer-kugawana/y2_pingdom_controller/workflows/ECR/badge.svg)](https://github.com/axel-springer-kugawana/y2_pingdom_controller/actions?query=workflow%3Aci+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/axel-springer-kugawana/y2_pingdom_controller)](https://goreportcard.com/report/github.com/axel-springer-kugawana/y2_pingdom_controller)

y2_pingdom_controller is a kubernetes controller that can create HTTP checks for [Pingdom](https://www.pingdom.com/).   
**Important:** The current version supports working with HTTP check only.

## Annotations

You can add these Kubernetes annotations to specific Ingress objects to customize their behavior.



|Name                       | Type | Optional | Example |
|---------------------------|------|------|------|
|pingdom.controller.yad2/apply|"true" or "false"| X | "true"
| pingdom.controller.yad2/resolution |string| X | "1"
|pingdom.controller.yad2/integrationids| string | √ | "92247"
|pingdom.controller.yad2/probe-filters| string | √ | "region: EU" 
|pingdom.controller.yad2/port| string | √ | "80"
|pingdom.controller.yad2/teamids| string | √ | "12345" 
|pingdom.controller.yad2/paused | string | √ | "false"
|pingdom.controller.yad2/verify-certificate |"true" or "false"| √ | "true"
