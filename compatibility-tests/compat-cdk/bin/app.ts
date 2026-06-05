#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { CloudStackTestStack } from '../lib/CloudStack-stack';

const app = new cdk.App();
new CloudStackTestStack(app, 'CloudStackTestStack', {
  env: {
    account: '000000000000',
    region: 'us-east-1',
  },
});
