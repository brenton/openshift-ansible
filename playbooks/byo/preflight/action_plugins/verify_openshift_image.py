from __future__ import (absolute_import, division, print_function)
__metaclass__ = type

import os

from ansible.plugins.action import ActionBase


class ActionModule(ActionBase):

    def run(self, tmp=None, task_vars=None):
        ''' handler verifying openshift images are available '''
        if task_vars is None:
            task_vars = dict()

        result = super(ActionModule, self).run(tmp, task_vars)

        component  = self._task.args.get('component', None)
        
        if component is None:
            result['failed'] = True
            result['msg'] = "component is required"
            return result

        docker_image_args = dict(
                name="openshift3/%s" % component,
            )

        result.update(self._execute_module(module_name='docker_image',
                module_args=docker_image_args, task_vars=task_vars))

        return result
