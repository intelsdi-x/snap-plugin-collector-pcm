
# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**




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

Namespace	|	Description
-----------	|	--------------
/intel/pcm/[Component]/EXEC	|	Instructions per nominal CPU cycle, i.e. in respect to the CPU frequency ignoring turbo and power saving
/intel/pcm/[Component]/IPC	|	Instructions per cycle. This measures how effectively you are using the core.
/intel/pcm/[Component]/FREQ	|	Frequency relative to nominal CPU frequency (“clockticks”/”invariant timer ticks”)
/intel/pcm/[Component]/AFREQ	|	Frequency relative to nominal CPU frequency excluding the time when the CPU is sleeping
/intel/pcm/[Component]/L3MISS	|	L3 cache line misses in millions
/intel/pcm/[Component]/L2MISS	|	L2 cache line misses in millions
/intel/pcm/[Component]/L3HIT	|	L3 Cache hit ratio (hits/reference)
/intel/pcm/[Component]/L2HIT	|	L2 Cache hit ratio (hits/reference)
/intel/pcm/[Component]/L3CLK	|	Very rough estimate of cycles lost to L3 cache misses vs. clockticks
/intel/pcm/[Component]/L2CLK	|	Very rough estimate of cycles lost to L2 cache misses vs. clockticks
/intel/pcm/[Component]/READ	|	Memory read traffic on this socket in GB
/intel/pcm/[Component]/WRITE	|	Memory write traffic on this socket in GB
/intel/pcm/[Component]/C[CoreNumber]res	|	Core residency
/intel/pcm/[Socket]/Proc_Energy_Joules	|	The energy consumed by the processor in Joules. Divide by the time to get the power consumption in watt
/intel/pcm/[Socket]/DRAM_Energy_Joules	|	The energy consumed by the DRAM attached to this socket in Joules. Divide by the time to get the power consumption in watt
/intel/pcm/[Socket]/TEMP	|	Thermal headroom in Kelvin (max design temperature – current temperature)
/intel/pcm/[System]/INST	|	Number of instructions retired
/intel/pcm/[System]/ACYC	|	Number of clockticks, This takes turbo and power saving modes into account.
/intel/pcm/[System]/TIME_ticks	|	Number of invariant clockticks. This is invariant to turbo and power saving modes.
/intel/pcm/[System]/PhysIPC	|	Instructions per cycle (IPC) multiplied by number of threads per core. See section "Core Cycles-per-Instruction (CPI) and Thread CPI" in Performance Insights to Intel® Hyper-Threading Technology for some background information.
/intel/pcm/[System]/PhysIPC%	|	Instructions per cycle (IPC) multiplied by number of threads per core relative to maximum IPC
/intel/pcm/[System]/INSTnom	|	Instructions per nominal cycle multiplied by number of threads per core
/intel/pcm/[System]/INSTnom%	|	Instructions per nominal cycle multiplied by number of threads per core relative to maximum IPC. The maximum IPC is 2 for Atom and 4 for all other supported processors.
/intel/pcm/[System]/TotalQPIin	|	QPI data traffic estimation (data traffic coming to CPU/socket through QPI links) in MB (1024*1024)
/intel/pcm/[System]/QPItoMC	|	Ratio of QPI traffic to memory traffic
/intel/pcm/[System]/TotalQPIout	|	QPI traffic estimation (data and non-data traffic outgoing from CPU/socket through QPI links) in MB (1024*1024)


Metrics exposed by "pcm" are system related and might be varied.

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
                "/intel/pcm/*": {}
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
ID: 44c01cd0-7133-49b1-a95c-a444db064b40
Name: Task-44c01cd0-7133-49b1-a95c-a444db064b40
State: Running
```

See sample output from `snaptel task watch <task_id>`

```
$ snaptel task watch 44c01cd0-7133-49b1-a95c-a444db064b40

Watching Task (44c01cd0-7133-49b1-a95c-a444db064b40):
NAMESPACE 						 DATA 		 TIMESTAMP
/intel/pcm/SKT0_Core_C-State/C0res% 			 1.13 		 2017-04-18 14:52:05.410848537 +0200 CEST
/intel/pcm/SKT0_Core_C-State/C1res% 			 26 		 2017-04-18 14:52:05.410835819 +0200 CEST
/intel/pcm/SKT0_Core_C-State/C3res% 			 0.00878 	 2017-04-18 14:52:05.410855565 +0200 CEST
/intel/pcm/SKT0_Core_C-State/C6res% 			 0.361 		 2017-04-18 14:52:05.410829547 +0200 CEST
/intel/pcm/SKT0_Core_C-State/C7res% 			 72.5 		 2017-04-18 14:52:05.410842087 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C10res% 		 0 		 2017-04-18 14:52:05.41061602 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C2res% 			 26.6 		 2017-04-18 14:52:05.410630665 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C3res% 			 0 		 2017-04-18 14:52:05.410623662 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C6res% 			 0 		 2017-04-18 14:52:05.410660845 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C7res% 			 0 		 2017-04-18 14:52:05.410637857 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C8res% 			 0 		 2017-04-18 14:52:05.410644744 +0200 CEST
/intel/pcm/SKT0_Package_C-State/C9res% 			 0 		 2017-04-18 14:52:05.410653227 +0200 CEST
/intel/pcm/Socket0/AFREQ 				 0.383 		 2017-04-18 14:52:05.410772622 +0200 CEST
/intel/pcm/Socket0/EXEC 				 0.005 		 2017-04-18 14:52:05.410818398 +0200 CEST
/intel/pcm/Socket0/FREQ 				 0.00432 	 2017-04-18 14:52:05.410797095 +0200 CEST
/intel/pcm/Socket0/IPC 					 1.16 		 2017-04-18 14:52:05.410767806 +0200 CEST
/intel/pcm/Socket0/L2HIT 				 0.57 		 2017-04-18 14:52:05.410782617 +0200 CEST
/intel/pcm/Socket0/L2MISS 				 0.472 		 2017-04-18 14:52:05.410802963 +0200 CEST
/intel/pcm/Socket0/L2MPI 				 0.00346 	 2017-04-18 14:52:05.410823254 +0200 CEST
/intel/pcm/Socket0/L3HIT 				 0.719 		 2017-04-18 14:52:05.410787351 +0200 CEST
/intel/pcm/Socket0/L3MISS 				 0.12 		 2017-04-18 14:52:05.410813474 +0200 CEST
/intel/pcm/Socket0/L3MPI 				 0.00088 	 2017-04-18 14:52:05.410763201 +0200 CEST
/intel/pcm/Socket0/READ 				 0.651 		 2017-04-18 14:52:05.410777707 +0200 CEST
/intel/pcm/Socket0/TEMP 				 68 		 2017-04-18 14:52:05.410807757 +0200 CEST
/intel/pcm/Socket0/WRITE 				 0.0908 	 2017-04-18 14:52:05.410792212 +0200 CEST
/intel/pcm/System/ACYC 					 118 		 2017-04-18 14:52:05.410506608 +0200 CEST
/intel/pcm/System/AFREQ 				 0.383 		 2017-04-18 14:52:05.410534795 +0200 CEST
/intel/pcm/System/EXEC 					 0.005 		 2017-04-18 14:52:05.410603378 +0200 CEST
/intel/pcm/System/FREQ 					 0.00432 	 2017-04-18 14:52:05.410567768 +0200 CEST
/intel/pcm/System/INST 					 136 		 2017-04-18 14:52:05.410524764 +0200 CEST
/intel/pcm/System/INSTnom 				 0.01 		 2017-04-18 14:52:05.410587925 +0200 CEST
/intel/pcm/System/INSTnom% 				 0.25 		 2017-04-18 14:52:05.410558073 +0200 CEST
/intel/pcm/System/IPC 					 1.16 		 2017-04-18 14:52:05.410501595 +0200 CEST
/intel/pcm/System/L2HIT 				 0.57 		 2017-04-18 14:52:05.410539947 +0200 CEST
/intel/pcm/System/L2MISS 				 0.472 		 2017-04-18 14:52:05.410598547 +0200 CEST
/intel/pcm/System/L2MPI 				 0.00346 	 2017-04-18 14:52:05.410578095 +0200 CEST
/intel/pcm/System/L3HIT 				 0.719 		 2017-04-18 14:52:05.410582892 +0200 CEST
/intel/pcm/System/L3MISS 				 0.12 		 2017-04-18 14:52:05.410496412 +0200 CEST
/intel/pcm/System/L3MPI 				 0.00088 	 2017-04-18 14:52:05.410593533 +0200 CEST
/intel/pcm/System/PhysIPC 				 2.32 		 2017-04-18 14:52:05.410563088 +0200 CEST
/intel/pcm/System/PhysIPC% 				 57.9 		 2017-04-18 14:52:05.410490594 +0200 CEST
/intel/pcm/System/READ 					 0.651 		 2017-04-18 14:52:05.410572346 +0200 CEST
/intel/pcm/System/TIME_ticks 				 3410 		 2017-04-18 14:52:05.410608886 +0200 CEST
/intel/pcm/System/WRITE 				 0.0908 	 2017-04-18 14:52:05.410529645 +0200 CEST
/intel/pcm/System_Core_C-States/C0res% 			 1.13 		 2017-04-18 14:52:05.410731258 +0200 CEST
/intel/pcm/System_Core_C-States/C1res% 			 26 		 2017-04-18 14:52:05.410723708 +0200 CEST
/intel/pcm/System_Core_C-States/C3res% 			 0.00878 	 2017-04-18 14:52:05.41075228 +0200 CEST
/intel/pcm/System_Core_C-States/C6res% 			 0.361 		 2017-04-18 14:52:05.410738707 +0200 CEST
/intel/pcm/System_Core_C-States/C7res% 			 72.5 		 2017-04-18 14:52:05.410745516 +0200 CEST
/intel/pcm/System_Pack_C-States/C10res% 		 0 		 2017-04-18 14:52:05.410696658 +0200 CEST
/intel/pcm/System_Pack_C-States/C2res% 			 26.6 		 2017-04-18 14:52:05.410716662 +0200 CEST
/intel/pcm/System_Pack_C-States/C3res% 			 0 		 2017-04-18 14:52:05.410667741 +0200 CEST
/intel/pcm/System_Pack_C-States/C6res% 			 0 		 2017-04-18 14:52:05.410674519 +0200 CEST
/intel/pcm/System_Pack_C-States/C7res% 			 0 		 2017-04-18 14:52:05.410703432 +0200 CEST
/intel/pcm/System_Pack_C-States/C8res% 			 0 		 2017-04-18 14:52:05.410689679 +0200 CEST
/intel/pcm/System_Pack_C-States/C9res% 			 0 		 2017-04-18 14:52:05.410709942 +0200 CEST
/intel/pcm/System_Pack_C-States/Proc_Energy_Joules 	 3.07 		 2017-04-18 14:52:05.410682756 +0200 CEST

```
(Keys `ctrl+c` terminate task watcher)

These data are published to file and stored there (in this example in /tmp/published_pcm).

Stop task:
```
$ snaptel task stop 44c01cd0-7133-49b1-a95c-a444db064b40
Task stopped:
ID: 44c01cd0-7133-49b1-a95c-a444db064b40
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

* Author: [Justin Guidroz](https://github.com/geauxvirtual)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
