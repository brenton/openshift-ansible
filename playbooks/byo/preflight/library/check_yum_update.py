#!/usr/bin/python

# tests whether a plain "yum update" will succeed

import os
import sys
import yum
from ansible.module_utils.basic import AnsibleModule


def main():
    module = AnsibleModule(
        argument_spec = dict()
    )
    def bail(error):
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

    updates = yb.update()
    # find out if there are any errors with the update
    # for error of type:
    #   1. dependency for update not found
    #    * retrieve the dependency tree
    #    * find the installed package(s) that required the missing dep
    #    * determine if any of these packages matter to openshift
    #    * build helpful error output
    #   2. conflicts among packages in available content
    #    * analyze dependency tree and build helpful error output
    #   3. other/unknown
    #    * just report the problem verbatim

    #thing = type(updates[0])
    #if updates:
    #    for pkg in updates:
    #        print "%s will be an update" % pkg
    module.exit_json(changed=False)

if __name__ == '__main__':
    main()
