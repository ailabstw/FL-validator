# What is FL validator

FL validator is a tool which can be used to validate whether your FL application container has been correctly implemented with our GRPC interface.

# Getting started

There are few things one will need to do before using the fl validator. The things are below.

1. Know how to make a valid fl application image. This can be learned from our Hello FL project. The following link will get you there.
      [Hello FL](https://gitlab.corp.ailabs.tw/federated-learning/hello-fl)

2. You will need to make a docker image (with all GRPC interfaces implemented) according to Hello FL. Let's call it **my_application**.

3. Alter the **docker-compose.yml** file in this project.

  * Change the image of the container ```app``` (which is ```registry.corp.ailabs.tw/federated-learning/hello-fl/edge:1.3.4``` now) to the image you have built.(at the 5th line in docker-compose.yml)

  * There are few environmental variables that need to be set.

    **LOCAL_MODEL_PATH** : is where you will place the model weight after you have trained a new weight per epoch.
  （The given example value of ```LOCAL_MODEL_PATH``` is ```/model/weight.ckpt```）

    **GLOBAL_MODEL_PATH** :is the path where you will load your pre-trained model or the globally merged model
  （The given example value of ```GLOBAL_MODEL_PATH``` is ```/model/merge.ckpt```）

  * **NVIDIA_VISIBLE_DEVICES** : need to be set as the index of the GPU card you will use.
  （The given example value of ```NVIDIA_VISIBLE_DEVICES``` is ```0```, because only choose first GPU card to do training.）

  * The docker-compose's mounting path（```volumes:```） also needs to be altered. There are two paths that need to be set . One is *model path*, and another is ```data path```.
    * **model path** is where you should put the merged global model weight and the local model weight. And this will be the folder path (without the model weight's name) of the ```LOCAL_MODEL_PATH```.
    * **data path** is where you should put your training datasets in correspond to the path you will load your datasets from.

* 4. **DRY_RUN** : set this to enable ```DRY_RUN``` or disable ```DRY_RUN```.
      There are two modes one can choose. One can set the value of ```DRY_RUN``` argument in the environment variables setting in the validator in our docker-compose.yml to decide whether to enable dry run or not.
        * ```DRY_RUN='True'``` enable DRY RUN mode. The validator will test the four GRPC interfaces of your ```3rd_application```.
        * ```DRY_RUN='False'``` disable DRY RUN and mode. The validator will test a full round of training with your ```3rd_application```.

      There is an important thing one needs to know. That is if one wants to open DRY RUN mode, their GRPC server implemented in ```3rd_application``` will need to handle ```DRY RUN mode message``` which is packaged in the context of every GRPC call validator makes.


      The ```DRY RUN mode message``` is a key-value pair (the key is a string and the value is bool) as below.
        ```bash
          "draftRun", true
        ```
        Or

        ```bash
          "draftRun", false
        ```
        All your GRPC interfaces should immediately return OK once you have received ```"draftRun", true``` after parsing the context contained in a GPPC call.


* 5. After you have set 1-3 above, you can simply run our validator with command as below.

```bash
docker-compose up -d
```


## What will be validated ?

  * **DataValidation**, **TrainInit**, **LocalTrain**, **TrainFinish** : this four GPPC interfaces will be validated. (There are five GRPC interfaces actually, but currently **TrainInterrupt** do not be used.) You will know whether your image has successfully implemented the basic interfaces to fit Ailab's federated learning system. A **report.json** file will be created at ```/var/reports/report.json``` (You can change this path in docker-compose.yml, the outside path correspond ```/var/report```)

<div align="left"><img src="./assets/validator_msc_2.png" style="width:100%"></img></div>

  * Whether your image can successfully do the first round of federated learning. If you have successfully done one round of federated learning, you will see the message below finally.
  And both your app and our validator will end up exit 0 soon.

<div align="left"><img src="./assets/validator_msc_sucess.png" style="width:100%"></img></div>


  * If your image have successfully implemented the log interface as the example (Hello FL), you will see a log file located at ```/var/logs/log.json```.

  <div align="left"><img src="./assets/validator_msc_1.png" style="width:100%"></img></div>
