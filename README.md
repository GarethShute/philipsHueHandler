# About philipsHueHandler
Hue handler acts as a simplified interface into your home Hue Lights network giving you control of the colour and brightness of the bulbs. 
Note: *You will need a hue home hub in order to make use of this application and this app will need to be running on a device connected to the same LAN as the hue home hub.*
## API endpoints
### /huelights/code (GET)
Use this endpoint to generate appKey for config file.
First call will generate the following reponse:
```
[  
    {  
        "error": {  
            "type": 101,  
            "address": "",  
            "description": "link button not pressed"  
        }  
    }  
]
```
You will then have to physically press the button on your HueHub to accept the request. Once you have pressed the button you can call this endpoint again and you should get a response in this format if successful:
```
[  
    {  
        "success":{  
            "username":<Code>,  
            "clientkey":<Code>  
        }  
    }  
]
```
You've now created an authorised *username*, which we'll use as the hue-application-key from now on. Place in the config file for the *appKey* value as descibed in the json config section below.
### /huelights/changecolour (POST) 
JSON body format:  
(X and Y values from CIE chromacity diagram)  
(mirek value from ~155 - 500 depending on warmth of light. Lower is cooler)  
Note: use mirek for white and xy for colour.  
```
{  
    "light" : <"colour" or "white">,                 -Required  
    "x" : 0.4605,                                    -Optional should be present if using Colour  
    "y" : 0.7,                                       -Optional should be present if using Colour  
    "mirek" : 200,                                   -Optional should be present if using White  
    "brightness" : 100                               -Required, percentage between 0-100  
}
```

## Json config format
*Note: You will need the grouped light id of the lights you wish to control*
```
{
    "hueIP" : <HUE_hub_IPAddress>,
    "appKey" : <APP_KEY>,
    "groupedLight" : <Grouped_light_ID>
    "port" : <port>                                  -Optional, will default to 8080 if not present
}
```
