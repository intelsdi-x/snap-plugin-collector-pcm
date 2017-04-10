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
* [golang 1.6+](https://golang.org/dl/)  (needed only for building)
* Root privileges (snapteld has to be running with root privileges for ability to collect data from PCM)
 
**Suggestions**
* To be able, to use PCM, [NMI watchdog](https://en.wikipedia.org/wiki/Non-maskable_interrupt) needs to be disabled. There are two ways to do this:
 * at running time: 
		`echo 0 > /proc/sys/kernel/nmi_watchdog`
 * or permanently: 
		`echo 'kernel.nmi_watchdog=0' >> /etc/sysctl.conf`
		
* Currently, Ubuntu 14.04 users have to manually compile PCM and add it to $PATH or export $SNAP_PCM_PATH to be able to use it.
* To be able to run PCM, access to CPUs MSRs (Model Specific Register) is needed. To obtain it, execute `modprobe msr` as root user


### Installation

#### To install Intel PCM:
Follow the instruction available at http://www.intel.com/software/pcm
To be sure that installed pcm.x works properly, try to execute `pcm.x /csv`

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

This builds the plugin in `./build/`

### Configuration and Usage

By default pcm executable binary are searched in the directories named by the PATH environment. 
Customize path to pcm executable is also possible by setting environment variable `export SNAP_PCM_PATH=/path/to/pcm/bin`

## Documentation

To learn more about Intel PCM visit http://www.intel.com/software/pcm

### Collected Metrics
This plugin has the ability to gather metrics for various components (like system, particular socket, dram etc.). Namespaces are constructed using following rule `/intel/pcm/[component name]/[metric name]`

Here are abbreviations for metric names [source](https://software.intel.com/en-us/blogs/2014/07/18/intel-pcm-column-names-decoder-ring):

The following metrics are available on all levels:

Metric	|	Explanation
-----------	|	--------------
EXEC	|	Instructions per nominal CPU cycle, i.e. in respect to the CPU frequency ignoring turbo and power saving
IPC	|	Instructions per cycle. This measures how effectively you are using the core.
FREQ	|	Frequency relative to nominal CPU frequency (“clockticks”/”invariant timer ticks”)
AFREQ	|	Frequency relative to nominal CPU frequency excluding the time when the CPU is sleeping
L3MISS	|	L3 cache line misses in millions
L2MISS	|	L2 cache line misses in millions
L3HIT	|	L3 Cache hit ratio (hits/reference)
L2HIT	|	L2 Cache hit ratio (hits/reference)
L3CLK	|	Very rough estimate of cycles lost to L3 cache misses vs. clockticks
L2CLK	|	Very rough estimate of cycles lost to L2 cache misses vs. clockticks

The following metrics are only available on socket and system level:

Metric	|	Explanation
-----------	|	--------------
READ	|	Memory read traffic on this socket in GB
WRITE	|	Memory read traffic on this socket in GB

The following metrics are only available on a socket level:

Metric	|	Explanation
-----------	|	--------------
Proc Energy (Joules)	|	The energy consumed by the processor in Joules. Divide by the time to get the power consumption in watt
DRAM Energy (Joules)	|	The energy consumed by the DRAM attached to this socket in Joules. Divide by the time to get the power consumption in watt
TEMP	|	Thermal headroom in Kelvin (max design temperature – current temperature)


The following metrics are only available on a system level:

Metric	|	Explanation
-----------	|	--------------
INST	|	Number of instructions retired
ACYC	|	Number of clockticks, This takes turbo and power saving modes into account.
TIME(ticks)	|	Number of invariant clockticks. This is invariant to turbo and power saving modes.
PhysIPC	|	Instructions per cycle (IPC) multiplied by number of threads per core. See section "Core Cycles-per-Instruction (CPI) and Thread CPI" in Performance Insights to Intel® Hyper-Threading Technology for some background information.
PhysIPC%	|	Instructions per cycle (IPC) multiplied by number of threads per core relative to maximum IPC
INSTnom	|	Instructions per nominal cycle multiplied by number of threads per core
INSTnom%	|	Instructions per nominal cycle multiplied by number of threads per core relative to maximum IPC. The maximum IPC is 2 for Atom and 4 for all other supported processors.
TotalQPIin	|	QPI data traffic estimation (data traffic coming to CPU/socket through QPI links) in MB (1024*1024)
QPItoMC	|	Ratio of QPI traffic to memory traffic
TotalQPIout	|	QPI traffic estimation (data and non-data traffic outgoing from CPU/socket through QPI links) in MB (1024*1024)

Example set of metric available on 2 socket platform:/intel/pcm/DRAM_Energy_Joules/SKT0
```
/intel/pcm/DRAM_Energy_Joules/SKT1
/intel/pcm/Proc_Energy_Joules/SKT0
/intel/pcm/Proc_Energy_Joules/SKT1
/intel/pcm/SKT0dataInSKT0dataIn_percent_SKT1dataInSKT1dataIn_percent_SKT0trafficOutSKT0trafficOut_percent_SKT1trafficOutSKT1trafficOut_percent_SKT0_Core_C-State/C0res%
/intel/pcm/SKT0dataInSKT0dataIn_percent_SKT1dataInSKT1dataIn_percent_SKT0trafficOutSKT0trafficOut_percent_SKT1trafficOutSKT1trafficOut_percent_SKT0_Core_C-State/C1res%
/intel/pcm/SKT0dataInSKT0dataIn_percent_SKT1dataInSKT1dataIn_percent_SKT0trafficOutSKT0trafficOut_percent_SKT1trafficOutSKT1trafficOut_percent_SKT0_Core_C-State/C3res%
/intel/pcm/SKT0dataInSKT0dataIn_percent_SKT1dataInSKT1dataIn_percent_SKT0trafficOutSKT0trafficOut_percent_SKT1trafficOutSKT1trafficOut_percent_SKT0_Core_C-State/C6res%
/intel/pcm/SKT0dataInSKT0dataIn_percent_SKT1dataInSKT1dataIn_percent_SKT0trafficOutSKT0trafficOut_percent_SKT1trafficOutSKT1trafficOut_percent_SKT0_Core_C-State/C7res%
/intel/pcm/SKT0_Package_C-State/C2res%
/intel/pcm/SKT0_Package_C-State/C3res%
/intel/pcm/SKT0_Package_C-State/C6res%
/intel/pcm/SKT0_Package_C-State/C7res%
/intel/pcm/SKT1_Core_C-State/C0res%
/intel/pcm/SKT1_Core_C-State/C1res%
/intel/pcm/SKT1_Core_C-State/C3res%
/intel/pcm/SKT1_Core_C-State/C6res%
/intel/pcm/SKT1_Core_C-State/C7res%
/intel/pcm/SKT1_Package_C-State/C2res%
/intel/pcm/SKT1_Package_C-State/C3res%
/intel/pcm/SKT1_Package_C-State/C6res%
/intel/pcm/SKT1_Package_C-State/C7res%
/intel/pcm/Socket0/AFREQ
/intel/pcm/Socket0/EXEC
/intel/pcm/Socket0/FREQ
/intel/pcm/Socket0/IPC
/intel/pcm/Socket0/L2HIT
/intel/pcm/Socket0/L2MISS
/intel/pcm/Socket0/L2MPI
/intel/pcm/Socket0/L3HIT
/intel/pcm/Socket0/L3MISS
/intel/pcm/Socket0/L3MPI
/intel/pcm/Socket0/READ
/intel/pcm/Socket0/TEMP
/intel/pcm/Socket0/WRITE
/intel/pcm/Socket1/AFREQ
/intel/pcm/Socket1/EXEC
/intel/pcm/Socket1/FREQ
/intel/pcm/Socket1/IPC
/intel/pcm/Socket1/L2HIT
/intel/pcm/Socket1/L2MISS
/intel/pcm/Socket1/L2MPI
/intel/pcm/Socket1/L3HIT
/intel/pcm/Socket1/L3MISS
/intel/pcm/Socket1/L3MPI
/intel/pcm/Socket1/READ
/intel/pcm/Socket1/TEMP
/intel/pcm/Socket1/WRITE
/intel/pcm/System/ACYC
/intel/pcm/System/AFREQ
/intel/pcm/System_Core_C-States/C0res%
/intel/pcm/System_Core_C-States/C1res%
/intel/pcm/System_Core_C-States/C3res%
/intel/pcm/System_Core_C-States/C6res%
/intel/pcm/System_Core_C-States/C7res%
/intel/pcm/System/EXEC
/intel/pcm/System/FREQ
/intel/pcm/System/INST
/intel/pcm/System/INSTnom
/intel/pcm/System/INSTnom%
/intel/pcm/System/IPC
/intel/pcm/System/L2HIT
/intel/pcm/System/L2MISS
/intel/pcm/System/L2MPI
/intel/pcm/System/L3HIT
/intel/pcm/System/L3MISS
/intel/pcm/System/L3MPI
/intel/pcm/System_Pack_C-States/C2res%
/intel/pcm/System_Pack_C-States/C3res%
/intel/pcm/System_Pack_C-States/C6res%
/intel/pcm/System_Pack_C-States/C7res%
/intel/pcm/System_Pack_C-States/DRAM_Energy_Joules
/intel/pcm/System_Pack_C-States/Proc_Energy_Joules
/intel/pcm/System/PhysIPC
/intel/pcm/System/PhysIPC%
/intel/pcm/System/QPItoMC
/intel/pcm/System/READ
/intel/pcm/System/TIME_ticks
/intel/pcm/System/TotalQPIin
/intel/pcm/System/TotalQPIout
/intel/pcm/System/WRITE
```
By default metrics are gathered once per second.

### Examples
Example running  pcm collector and writing data to file. Notice that snapteld has to be running with root privileges, for ability to collect data from PCM

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `sudo service snap-telemetry start`
* systemd: `sudo systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`

Download and load snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-pcm/latest/linux/x86_64/snap-plugin-collector-pcm
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-pcm
$ snaptel plugin load snap-plugin-publisher-file
```

See available metrics for your system:
```
$ snaptel metric list
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
                "/intel/pcm/System/IPC": {},
                "/intel/pcm/System/L2HIT": {},
                "/intel/pcm/System/L2MISS": {},
                "/intel/pcm/System/EXEC": {},
                "/intel/pcm/System/FREQ": {},
                "/intel/pcm/System/INST": {},
                "/intel/pcm/System/INSTnom": {},
                "/intel/pcm/System/INSTnom%": {},
                "/intel/pcm/System/L3HIT": {},
                "/intel/pcm/System/L3MISS": {},
                "/intel/pcm/System/PhysIPC": {},
                "/intel/pcm/System/PhysIPC%": {},
                "/intel/pcm/System/READ": {},
                "/intel/pcm/Socket0/TEMP": {},
                "/intel/pcm/System/TIME_ticks": {},
                "/intel/pcm/System/WRITE": {}
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
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
snaptel task create -t examples/tasks/pcm-file.json
Using task manifest to create task
Task created
ID: 156366f2-e497-4c10-ad22-560fc71986af
Name: Task-156366f2-e497-4c10-ad22-560fc71986af
State: Running
```

See sample output from `snaptel task watch <task_id>`

```
$ snaptel task watch 156366f2-e497-4c10-ad22-560fc71986af

Watching Task (156366f2-e497-4c10-ad22-560fc71986af):
NAMESPACE                        DATA            TIMESTAMP
^Cntel/pcm/Socket0/TEMP          54              2017-04-10 15:45:49.983616028 +0200 CEST
/intel/pcm/System/EXEC           0.0156          2017-04-10 15:45:49.983613192 +0200 CEST
/intel/pcm/System/FREQ           0.0173          2017-04-10 15:45:49.983613905 +0200 CEST
/intel/pcm/System/INST           323             2017-04-10 15:45:49.983604288 +0200 CEST
/intel/pcm/System/INSTnom        0.0311          2017-04-10 15:45:49.983609734 +0200 CEST
/intel/pcm/System/INSTnom%       0.778           2017-04-10 15:45:49.983605404 +0200 CEST
/intel/pcm/System/IPC            0.9             2017-04-10 15:45:49.983601543 +0200 CEST
/intel/pcm/System/L2HIT          0.454           2017-04-10 15:45:49.983610516 +0200 CEST
/intel/pcm/System/L2MISS         1.46            2017-04-10 15:45:49.983606359 +0200 CEST
/intel/pcm/System/L3HIT          0.767           2017-04-10 15:45:49.983603143 +0200 CEST
/intel/pcm/System/L3MISS         0.323           2017-04-10 15:45:49.983616942 +0200 CEST
/intel/pcm/System/PhysIPC        1.8             2017-04-10 15:45:49.983611261 +0200 CEST
/intel/pcm/System/PhysIPC%       45              2017-04-10 15:45:49.983612275 +0200 CEST
/intel/pcm/System/READ           0.742           2017-04-10 15:45:49.983614903 +0200 CEST
/intel/pcm/System/TIME_ticks     2600            2017-04-10 15:45:49.983607786 +0200 CEST
/intel/pcm/System/WRITE          0.0872          2017-04-10 15:45:49.983608728 +0200 CEST

```
(Keys `ctrl+c` terminate task watcher)

These data are published to file and stored there (in this example in /tmp/published_pcm).

Stop task:
```
$ snaptel task stop 156366f2-e497-4c10-ad22-560fc71986af
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
