# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins
import copy
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from . import _utilities

__all__ = ['ProjectMemberArgs', 'ProjectMember']

@pulumi.input_type
class ProjectMemberArgs:
    def __init__(__self__, *,
                 member_id: pulumi.Input[builtins.str],
                 project_id: pulumi.Input[builtins.str]):
        """
        The set of arguments for constructing a ProjectMember resource.
        """
        pulumi.set(__self__, "member_id", member_id)
        pulumi.set(__self__, "project_id", project_id)

    @property
    @pulumi.getter
    def member_id(self) -> pulumi.Input[builtins.str]:
        return pulumi.get(self, "member_id")

    @member_id.setter
    def member_id(self, value: pulumi.Input[builtins.str]):
        pulumi.set(self, "member_id", value)

    @property
    @pulumi.getter
    def project_id(self) -> pulumi.Input[builtins.str]:
        return pulumi.get(self, "project_id")

    @project_id.setter
    def project_id(self, value: pulumi.Input[builtins.str]):
        pulumi.set(self, "project_id", value)


@pulumi.type_token("cloudbeaver:index:ProjectMember")
class ProjectMember(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 member_id: Optional[pulumi.Input[builtins.str]] = None,
                 project_id: Optional[pulumi.Input[builtins.str]] = None,
                 __props__=None):
        """
        Create a ProjectMember resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: ProjectMemberArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a ProjectMember resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param ProjectMemberArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ProjectMemberArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 member_id: Optional[pulumi.Input[builtins.str]] = None,
                 project_id: Optional[pulumi.Input[builtins.str]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ProjectMemberArgs.__new__(ProjectMemberArgs)

            if member_id is None and not opts.urn:
                raise TypeError("Missing required property 'member_id'")
            __props__.__dict__["member_id"] = member_id
            if project_id is None and not opts.urn:
                raise TypeError("Missing required property 'project_id'")
            __props__.__dict__["project_id"] = project_id
        super(ProjectMember, __self__).__init__(
            'cloudbeaver:index:ProjectMember',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'ProjectMember':
        """
        Get an existing ProjectMember resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = ProjectMemberArgs.__new__(ProjectMemberArgs)

        __props__.__dict__["member_id"] = None
        __props__.__dict__["project_id"] = None
        return ProjectMember(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter
    def member_id(self) -> pulumi.Output[builtins.str]:
        return pulumi.get(self, "member_id")

    @property
    @pulumi.getter
    def project_id(self) -> pulumi.Output[builtins.str]:
        return pulumi.get(self, "project_id")

