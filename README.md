# Applifier

[![Build Status](https://cloud.drone.io/api/badges/elahe-dastan/applifier/status.svg)](https://cloud.drone.io/elahe-dastan/applifier)

## Protocol
Server and client can talk to each other using the following set of rules:<br/>

Colons can be used to align columns.

| Client              | Server           | Example                                                          |
| ------------------- |:-----------------| -----------------------------------------------------------------|
| WhoAmI              | WhoAmI,id        | client:"WhoAmI"<br/>server:"WhoAmI,1                             |
| List                | List,id1-id2,... | client:"List"<br/>server:"List,2-3-4                             |
| Send,id-id-...,body | Send,body        | client:"Send,2-3,I will be late<br/>server:"Send,I will be late" |
