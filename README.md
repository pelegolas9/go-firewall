# go-firewall
A basic CLI that emulates the behavior of a Firewall

# About the project

This project aims to better understand Go and network systems. The main goals are:
- Learn to manipulate network packets in Go
- Explore baisc networking concepts
- Apply the knowledge to a real-world scenario and something useful

# How to use it

When you run the main files, the CLI starts to listen on a specified HOST and PORT (hard-coded as localhost and port 3333).
Type addrule to add an incoming rule to block/allow netork packets from the origin.
Use the send_packets.go script to customize packets you send.
The CLI response is
packet allowed: {...} or
packet blocked: {...}

# Prerequisistes

You need Go installed (version >=1.20)
Admin/root permissions to manpulate network packets

# Installation and usage

Clone the repo
Navigate to go-core directory
Run the two main files with: go run .
