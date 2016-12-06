#!/usr/bin/python

import os
import sys
import yum
from ansible.module_utils.basic import AnsibleModule

"""
An ansible module for determining if more than one minor version
of any atomic-openshift package is available, which would indicate
that multiple repos are enabled for different versions of the same
thing which may cause problems.

Also, determine if the version requested is available down to the
precision requested.
"""


def main():
    module = AnsibleModule(
        argument_spec=dict(
            version=dict(required=True)
        ),
        supports_check_mode=True
    )
    sys.stdout = os.devnull  # mute yum so it doesn't break our output
    # sys.stderr = os.devnull # mute yum so it doesn't break our output

    def _unmute():
        sys.stdout = sys.__stdout__

    def bail(error):
        _unmute()
        module.fail_json(msg=error)

    yb = yum.YumBase()

    # search for package versions available for aos pkgs
    expected_pkgs = ["atomic-openshift",
                     "atomic-openshift-node", "atomic-openshift-master"]
    try:
        pkgs = yb.pkgSack.returnPackages(patterns=expected_pkgs)
    except yum.Errors.PackageSackError as e:
        # you only hit this if *none* of the packages are available
        bail("Unable to find any atomic-openshift packages. \nCheck your subscription and repo settings. \n%s" % e)

    # determine what level of precision we're expecting for the version
    expected_version = module.params['version']
    if expected_version.startswith('v'):  # v3.3 => 3.3
        expected_version = expected_version[1:]
    numDots = expected_version.count('.')

    pkgsByNameVersion = {}
    pkgsPreciseVersionFound = {}
    for pkg in pkgs:
        # get expected version precision
        match_version = '.'.join(pkg.version.split('.')[:numDots + 1])
        if match_version == expected_version:
            pkgsPreciseVersionFound[pkg.name] = True
        minor_version = '.'.join(pkg.version.split(
            '.')[:2])  # get x.y version precision
        if pkg.name not in pkgsByNameVersion:
            pkgsByNameVersion[pkg.name] = {}
        pkgsByNameVersion[pkg.name][minor_version] = True

    # see if any packages couldn't be found at requested version
    # see if any packages are available in more than one minor version
    not_found = []
    multi_found = []
    for name in expected_pkgs:
        if name not in pkgsPreciseVersionFound:
            not_found.append(name)
        if name in pkgsByNameVersion and len(pkgsByNameVersion[name]) > 1:
            multi_found.append(name)
    if not_found:
        msg = "Not all of the required packages are available at requested version %s:\n" % expected_version
        for name in not_found:
            msg += "  %s\n" % name
        bail(msg + "Please check your subscriptions and enabled repositories.")
    if multi_found:
        msg = "Multiple minor versions of these packages are available\n"
        for name in multi_found:
            msg += "  %s\n" % name
        bail(msg + "There should only be one OpenShift version's repository enabled at a time.")

    _unmute()
    module.exit_json(changed=False)


if __name__ == '__main__':
    main()
