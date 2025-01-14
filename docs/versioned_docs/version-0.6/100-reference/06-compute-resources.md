---
title: Compute resources
---

Acorn comes with a number of features that allow you define how your application workloads will be scheduled at deploy-time.

## Memory
You can configure Acorn apps to have a set memory upon startup.

This allows you to specify memory that an Acorn will abide by when being created. There are a number of ways to set this so that you have full control over defaults, workloads, and maximums. In order of precedence, the ways to set memory are when you:

1. [Run an Acorn](50-running/55-compute-resources.md)
2. [Author an Acornfile](38-authoring/03-containers.md#memory)
3. [Install Acorn](30-installation/02-options.md#memory)

:::note
When installing Acorn, you can also specify `--workload-memory-maximum`. This flag sets a maximum that when exceeded prevents the offending Acorn from being installed.
:::

### Valid memory values
Supported value formats for memory flags include
- 1_234 ->`1234`
- 5M -> `5_000_000`
- 1.5Gi ->`1_610_612_736`
- 0x1000_0000 -> `268_435_456`

These all translate into an exact amount of bytes. We encourage you use the binary representation of large quantities of bytes when interacting with memory such as `Ki`, `Mi`, `Gi`, and `Pi`.

#### No limit
It is possible to set any of these flags to have no limit on memory by simply setting their value to `0`. However, running an Acorn with its memory set to `0` while the `--workload-memory-maximum` is not set to `0` will roughly translate to "use as much memory as allowed". For example, say that we install Acorn with a non-zero `--workload-memory-maximum`.

```console
acorn install --workload-memory-maximum 512Mi
```

Then we try to run an Acorn with its memory set to 0.

```console
acorn run --memory 0 foo
```

When the `foo` Acorn gets provisioned, all of its containers will have their memory set to `512Mi` (the `--workload-memory-maximum` we set prior).

:::note
This same interaction will occur if the `--workload-memory-default` is set to 0 (which it is by default)
:::

## Compute Classes
You can configure Acorn apps to have a set compute class upon startup.

Setting a compute class allows you to define what the infrastructure providing your Acorn workloads should look like. Things that compute classes control include:

- What OS/Architecure your workloads will run on
- How much memory is minimal, maximal, default and allowed
- How many vCPUs should be allocated

:::info
You are not able to set vCPUs directly. This is an intentional abstraction and instead vCPUs are calculated based on of the amount of memory specified for a workload.
:::

### Using a Compute Class
You can see the compute classes available in your current project using the CLI. 

```console
$ acorn offerings computeclasses
NAME          DEFAULT   MEMORY RANGE      MEMORT DEFAULT   DESCRIPTION         
default       *         512Mi-1Gi         1Gi              Default ComputeClass
non-default             0-1Gi             512Mi            Non-default ComputeClass
unrestricted            Unrestricted      512Mi            Unrestricted ComputeClass
specific                128Mi,512Mi,1Gi   128Mi            Specific ComputeClass
```

Breaking this down:

- `DEFAULT` indicates whether a ComputeClass is the default and will be used if no ComputeClass is defined.
- `MEMORY_DEFAULT` indicates the memory that will be requested for a workload if none is specified.
- `MEMORY_RANGE` indicates what memory values are available to use, and can take two forms.
    1. `-` denotes a range and any value within it can be specified.
    2. `,` denotes specific values that are the only ones allowed.

Specifying compute classes can be done in the Acornfile (using the [class property](03-acornfile.md#class) for containers) or at runtime (using the [--compute-class flag](50-running/55-compute-resources.md#compute-class)). If you do not specify a compute class, the default compute class for the project will be used. If there is no default for the project, the default for the cluster will be used. Finally, if there is no cluster default then no compute class will be used. Depending on the compute class that is used, the memory that you specify may be in contention with its requirements. Should that happen Acorn will provide a descriptive error message to ammend any issues.

:::note
Looking to manage a compute class? This should only be done if you are (or are in communication with) an administrator of Acorn. You can read more information about managing compute classes [here](02-admin/03-computeclasses.md)
:::