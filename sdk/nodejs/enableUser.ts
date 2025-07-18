// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

export function enableUser(args: EnableUserArgs, opts?: pulumi.InvokeOptions): Promise<EnableUserResult> {
    opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts || {});
    return pulumi.runtime.invoke("cloudbeaver:index:enableUser", {
        "enabled": args.enabled,
        "user_name": args.user_name,
    }, opts);
}

export interface EnableUserArgs {
    enabled: boolean;
    user_name: string;
}

export interface EnableUserResult {
    readonly enabled: boolean;
    readonly user_name: string;
}
export function enableUserOutput(args: EnableUserOutputArgs, opts?: pulumi.InvokeOutputOptions): pulumi.Output<EnableUserResult> {
    opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts || {});
    return pulumi.runtime.invokeOutput("cloudbeaver:index:enableUser", {
        "enabled": args.enabled,
        "user_name": args.user_name,
    }, opts);
}

export interface EnableUserOutputArgs {
    enabled: pulumi.Input<boolean>;
    user_name: pulumi.Input<string>;
}
