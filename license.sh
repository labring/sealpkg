#!/bin/bash
go install github.com/google/addlicense@latest
addlicense -y 2023 -c "sealos." -f template/LICENSE pkg cmd types
