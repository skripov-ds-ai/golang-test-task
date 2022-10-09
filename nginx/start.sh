#!/bin/bash

service nginx start & consul-template -config="/scripts/consul-template-config.hcl"
