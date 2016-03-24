# snap collector plugin - pcm

This plugin collects  metrics from PCM (Intel Performance Counter Monitor)

It is used in the [snap framework] (http://github.com/intelsdi-x/snap).


1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

In order to use this plugin user is required to have PCM installed in system.

### System Requirements

* [Intel PCM] (http://www.intel.com/software/pcm)
* [golang 1.5+](https://golang.org/dl/)
* Root privileges (snapd has to be running with root privileges for ability to collect data from PCM)
 
**Suggestions**
* To be able, to use PCM, [NMI watchdog](https://en.wikipedia.org/wiki/Non-maskable_interrupt) needs to be disabled. There are two ways to do this:
 * at running time: 
		`echo 0 > /proc/sys/kernel/nmi_watchdog`
 * or permanently: 
		`echo 'kernel.nmi_watchdog=0' >> /etc/sysctl.conf`
		
* Currently, Ubuntu 14.04 users have to manually compile PCM and add it to $PATH or export $SNAP_PCM_PATH to be able to use it.

### Installation

#### To install Intel PCM:
Follow the instruction available at http://www.intel.com/software/pcm

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-pcm  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-pcm.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```

This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

By default pcm executable binary are searched in the directories named by the PATH environment. 
Customize path to pcm executable is also possible by setting environment variable `export SNAP_PCM_PATH=/path/to/pcm/bin`

## Documentation

To learn more about Intel PCM visit http://www.intel.com/software/pcm

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description 
------------ | -------------
/intel/pcm/ACYC|  Number of clockticks. This takes turbo and power saving modes into account.    
/intel/pcm/AFREQ| Frequency relative to nominal CPU frequency excluding the time when the CPU is sleeping.
/intel/pcm/C0res%| Core 0 residency
/intel/pcm/C1res%| Core 1 residency
/intel/pcm/C2res%| Core 2 residency
/intel/pcm/C3res%| Core 3 residency
/intel/pcm/Cres%| Cores residency
/intel/pcm/EXEC| Instructions per nominal CPU cycle, i.e. in respect to the CPU frequency ignoring turbo and power saving.
/intel/pcm/FREQ| Frequency relative to nominal CPU frequency, equals clockticks/invariant_timer_ticks.          
/intel/pcm/INST| Number of instructions retired            
/intel/pcm/INSTnom| Instructions per nominal cycle multiplied by number of threads per core.              
/intel/pcm/INSTnom%| Instructions per nominal cycle multiplied by number of threads per core relative to maximum IPC. The maximum IPC is 2 for Atom and 4 for all other supported processors.       
/intel/pcm/IPC| Instructions per cycle, this measures how effectively you are using the core.        
/intel/pcm/L2CLK| Very rough estimate of cycles lost to L2 cache misses vs. clockticks.
/intel/pcm/L2HIT| L2 cache hit ratio (0.00-1.00)            
/intel/pcm/L2MISS| L2 cache line misses
/intel/pcm/L2MPI| L2 cache misses per instruction    
/intel/pcm/L3CLK| Very rough estimate of cycles lost to L3 cache misses vs. clockticks.   
/intel/pcm/L3HIT| L3 cache hit ratio (0.00-1.00)          
/intel/pcm/L3MISS| L3 cache line misses        
/intel/pcm/L3MPI|  L3 cache misses per instruction            
/intel/pcm/PhysIPC| Instructions per cycle (IPC) multiplied by number of threads per core.    
/intel/pcm/PhysIPC%| Instructions per cycle (IPC) multiplied by number of threads per core relative to maximum IPC.          
/intel/pcm/ProcEnergy(Joules)| The energy consumed by the processor in Joules. Divide by the time to get the power consumption in watt
/intel/pcm/READ| Bytes read from memory controller in GBytes.
/intel/pcm/SKT0| CPU energy in Joules per socket 0
/intel/pcm/TEMP| Temperature reading in degree Celsius relative to the TjMax temperature (thermal headroom; max_design_temp - current_temp)
/intel/pcm/TIME(ticks)| Number of invariant clockticks. This is invariant to turbo and power saving modes.
/intel/pcm/WRITE| Bytes written to memory controller in GBytes.

Metrics exposed by "pcm" are system related and might be varied.

By default metrics are gathered once per second.

### Examples
Example running  pcm collector and writing data to file. Notice that snapd has to be running with root privileges, for ability to collect data from PCM

In one terminal window, open the snap daemon:
```
$ snapd -l 1 -t 0
```

In another terminal window, load pcm plugin for collecting:
```
$ snapctl plugin load $SNAP_PCM_PLUGIN_DIR/build/rootfs/snap-plugin-collector-pcm
Plugin loaded
Name: pcm
Version: 6
Type: collector
Signed: false
Loaded Time: Wed, 02 Dec 2015 07:57:33 EST
```
See available metrics for your system:
```
$ snapctl metric list
```

Load file plugin for publishing:
```
$ snapctl plugin load $SNAP_DIR/build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Wed, 02 Dec 2015 07:58:47 EST
```

Create a task JSON file (exemplary file in examples/tasks/pcm-file.json):  
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/pcm/IPC": {},
                "/intel/pcm/L2HIT": {},
                "/intel/pcm/L2MISS": {},
                "/intel/pcm/EXEC": {},
                "/intel/pcm/FREQ": {},
                "/intel/pcm/INST": {},
                "/intel/pcm/INSTnom": {},
                "/intel/pcm/INSTnom%": {},
                "/intel/pcm/L3HIT": {},
                "/intel/pcm/L3MISS": {},
                "/intel/pcm/PhysIPC": {},
                "/intel/pcm/PhysIPC%": {},
                "/intel/pcm/ProcEnergy(Joules)": {},
                "/intel/pcm/READ": {},
                "/intel/pcm/SKT0": {},
                "/intel/pcm/TEMP": {},
                "/intel/pcm/TIME(ticks)": {},
                "/intel/pcm/WRITE": {}
            },
            "config": {
                "/intel/pcm": {
                    "user": "root",
                    "password": "secret"
                }
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "plugin_version": 3,
                    "config": {
                        "file": "/tmp/published_pcm"
                    }
                }
            ]
        }
    }
}
```

Create a task:
```
snapctl task create -t $SNAP_PCM_PLUGIN_DIR/examples/tasks/pcm-file.json
Using task manifest to create task
Task created
ID: 156366f2-e497-4c10-ad22-560fc71986af
Name: Task-156366f2-e497-4c10-ad22-560fc71986af
State: Running
```

See sample output from `snapctl task watch <task_id>`

```
$ snapctl task watch 156366f2-e497-4c10-ad22-560fc71986af

Watching Task (156366f2-e497-4c10-ad22-560fc71986af):
NAMESPACE                        DATA            TIMESTAMP                                       SOURCE
/intel/pcm/EXEC                  0.0138          2015-12-02 08:19:46.001151927 -0500 EST         gklab-108-166
/intel/pcm/FREQ                  0.00639         2015-12-02 08:19:46.001150464 -0500 EST         gklab-108-166
/intel/pcm/INST                  379             2015-12-02 08:19:46.001150975 -0500 EST         gklab-108-166
/intel/pcm/INSTnom               0.0276          2015-12-02 08:19:46.001147704 -0500 EST         gklab-108-166
/intel/pcm/INSTnom%              0.691           2015-12-02 08:19:46.001148234 -0500 EST         gklab-108-166
/intel/pcm/IPC                   2.16            2015-12-02 08:19:46.001148772 -0500 EST         gklab-108-166
/intel/pcm/L2HIT                 0.483           2015-12-02 08:19:46.00114933 -0500 EST          gklab-108-166
/intel/pcm/L2MISS                0.719           2015-12-02 08:19:46.001151493 -0500 EST         gklab-108-166
/intel/pcm/L3HIT                 0.423           2015-12-02 08:19:46.001152449 -0500 EST         gklab-108-166
/intel/pcm/L3MISS                0.415           2015-12-02 08:19:46.001144495 -0500 EST         gklab-108-166
/intel/pcm/PhysIPC               4.33            2015-12-02 08:19:46.001145292 -0500 EST         gklab-108-166
/intel/pcm/PhysIPC%              108             2015-12-02 08:19:46.001149828 -0500 EST         gklab-108-166
/intel/pcm/ProcEnergy(Joules)    8.46            2015-12-02 08:19:46.001145857 -0500 EST         gklab-108-166
/intel/pcm/READ                  0.084           2015-12-02 08:19:46.00114662 -0500 EST          gklab-108-166
/intel/pcm/SKT0                  8.46            2015-12-02 08:19:46.001152938 -0500 EST         gklab-108-166
/intel/pcm/TEMP                  70              2015-12-02 08:19:46.001153401 -0500 EST         gklab-108-166
/intel/pcm/TIME(ticks)           3430            2015-12-02 08:19:46.001153955 -0500 EST         gklab-108-166
/intel/pcm/WRITE                 0.0563          2015-12-02 08:19:46.00114718 -0500 EST          gklab-108-166
```
(Keys `ctrl+c` terminate task watcher)

These data are published to file and stored there (in this example in /tmp/published_pcm).

Stop task:
```
$ $SNAP_PATH/bin/snapctl task stop 156366f2-e497-4c10-ad22-560fc71986af
Task stopped:
ID: 156366f2-e497-4c10-ad22-560fc71986af
```

### Roadmap
This plugin is in active development. As we launch this plugin, we have a few items in mind for the next release:
- [ ] Use channels instead "for" loop to execute pcm cmd

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-pcm/issues) 
and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-pcm/pulls).

## Community Support
This repository is one of **many** plugins in the **snap**, a powerful telemetry agent framework. See the full project at 
http://github.com/intelsdi-x/snap. To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).


## Contributing
We love contributions! :heart_eyes:

There is more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).


## Acknowledgements

* Author: [Izabella Raulin](https://github.com/IzabellaRaulin)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.