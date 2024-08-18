# tito-client
Simple Golang client for ti.to REST APIs.

This is my take on a Golang 'hello world'. It creates bulk ticket discount codes for events in [ti.to](https://ti.to/home).

This is very much a work in progress and is specific to my particular circumstances - to create speaker discount codes for the two events I help to organise: DevOpsDays London and Fast Flow Conf. 

Most values are currently hard-coded; the API key is passed as an argument.
The names of the speakers are supplied in a text file.

It's currently very basic but it meant that I didn't have to manually create over 120 discount codes and was a good excuse to have a play with Go.
