# FL validator

requirement:
* go 1.15

what will be validated ?
* Three GPPC interface : TrainInit, LocalTrain, TrainFinish will be validated.

setting:
* FL application 's GPRC server should listen: 0.0.0.0:7878
* FL application 's GPRC client should connect: 0.0.0.0:8787

usage:

```bash
go run .
```

validation sucessful message:

```plainText
"All FL validation completed . Congrats. "
```