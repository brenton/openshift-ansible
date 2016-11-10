#!/usr/bin/python

"""
Ansible module to test whether a yum update or install will succeed,
without actually performing it or running yum.
parameters:
  packages: (optional) A list of package names to install or update.
            If omitted, all installed RPMs are considered for updates.
"""

import os
import sys
import yum
from ansible.module_utils.basic import AnsibleModule


def main():
    module = AnsibleModule(
        argument_spec = dict(
            packages  = dict(type='list', default=[])
        ),
        supports_check_mode = True
    )
    sys.stdout = os.devnull # mute yum so it doesn't break our output

    def _unmute():
        sys.stdout = sys.__stdout__

    def bail(error):
        _unmute()
        module.fail_json(msg=error)

    yb = yum.YumBase()
    # determine if the existing yum configuration is valid
    try:
        yb.repos.populateSack(mdtype='metadata', cacheonly=1)
    # for error of type:
    #   1. can't reach the repo URL(s)
    except yum.Errors.NoMoreMirrorsRepoError as e:
        bail("Error getting data from at least one yum repository: %s" % e)
    #   2. invalid repo definition
    except yum.Errors.RepoError as e:
        bail("Error with yum repository configuration: %s" % e)
    #   3. other/unknown
    #    * just report the problem verbatim
    except:
        bail("Unexpected error with yum repository: %s" % sys.exc_info()[1])

    packages = module.params['packages']
    noSuchPkg = []
    for pkg in packages:
        try:
            yb.install(name=pkg)
        except yum.Errors.InstallError as e:
            noSuchPkg.append(pkg)
        except:
            bail("Unexpected error with yum install/update: %s" % sys.exc_info()[1])
    if not packages:
        # no packages requested means test a yum update of everything
        yb.update()
    elif noSuchPkg:
        # wanted specific packages to install but some aren't available
        userMsg = "Cannot install all of the necessary packages. Unavailable:\n"
        for pkg in noSuchPkg:
            userMsg += "  %s\n" % pkg
        userMsg += "You may need to enable one or more repos to make this content available."
        bail(userMsg)

    try:
        txnResult, txnMsgs = yb.buildTransaction()
    except:
        bail("Unexpected error during dependency resolution: %s" % sys.exc_info()[1])

    # find out if there are any errors with the update/install
    if txnResult == 0: # "normal exit" meaning there's nothing to install/update
        pass
    elif txnResult == 1: # error with transaction
        userMsg = "Could not perform yum update.\n"
        if len(txnMsgs) > 0:
            userMsg += "Errors from resolution:\n"
            for msg in txnMsgs:
                userMsg += "  %s\n" % msg
        bail(userMsg)
    # TODO: it would be nice depending on the problem:
    #   1. dependency for update not found
    #    * construct the dependency tree
    #    * find the installed package(s) that required the missing dep
    #    * determine if any of these packages matter to openshift
    #    * build helpful error output
    #   2. conflicts among packages in available content
    #    * analyze dependency tree and build helpful error output
    #   3. other/unknown
    #    * report the problem verbatim
    #    * add to this list as we come across problems we can clearly diagnose
    elif txnResult == 2: # everything resolved fine
        pass
    else:
        bail("Unknown error(s) from dependency resolution. Exit Code: %d:\n%s" % (txnResult, txnMsgs))

    _unmute()
    module.exit_json(changed=False)

if __name__ == '__main__':
    main()
