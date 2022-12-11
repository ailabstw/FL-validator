# FL validator

FL validator is a tool which can be used to validate whether your container is valid or not.

# Getting started

There are few things one will need to do before using fl validator as below.

* 1. Know how to make a valid fl application image. This can learn from our Hello FL project. The following link will get you to there.

* 2. If you have implemented all interface in your code, then package it as a docker image.

* 3. Alter the *docker-compose.yml* in this project.
  * Change image name of app to the container name you just built from your Dockerfile (usually at the 5th line)
  * There are 2 enviroments variables *LOCAL_MODEL_PATH* and *GLOBAL_MODEL_PATH* need to be set. LOCAL_MODEL_PATH is where you will place the model weight after you have trained a new weight per epoch. And GLOBAL_MODEL_PATH is the path where you will load your pre-trained model or the globally merged model（the epochs after first epoch）.
  * *NVIDIA_VISIBLE_DEVICES* need to be set as the index of GPU card you will use.

  * The mounting path also need to be altered. There are two path need to be set . One is *model path*, and another is *data path*.
    * *model path* is where you should put the merged global model weight and the local model weight.


* 4. After you have done 1-3 above, you can simply run our validator with command as below.

```bash
docker-compose up -d
```

## What will be validated ?
  * Four GPPC interface : DataValidation, TrainInit, LocalTrain, TrainFinish will be validated. You will know whether your image have sucessfully implemented the basic interface to fit our federated learning system. A report.json file will be created at */var/reports/report.json*

<div align="left"><img src="./assets/validator_msc_2.png" style="width:100%"></img></div>

  * Whether your image can sucessfully do the first round of federated learning. If you have sucessfully done one round of federated learning, you will see the message as below finally.
  And both your app and our validator will end up exit 0 soon.

<div align="left"><img src="./assets/validator_msc_sucess.png" style="width:100%"></img></div>


  * If your image have sucessfully implemented the log interface as the example (Hello FL), you will see a log file located at */var/logs/log.json*.

  <div align="left"><img src="./assets/validator_msc_1.png" style="width:100%"></img></div>