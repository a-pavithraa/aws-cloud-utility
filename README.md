# aws-cloud-utility

 It is always advisable to use IAC for spinning up resources  but sometimes for ad-hoc testing , I use AWS console for creating resources.  There was one particular instance where I had kept my LB in ap-south-1 for couple of months incurring huge bill. I tend to check in us-east-region and skip other regions. I am learning Go and thought to implement a small utility for checking whether I have any running ec2 or load balancer in the 3-4 regions I mainly create resources. 

 Regions, port are configured in config.json

```json
{
  "port": ":8080",
  "logLevel": "INFO",
  "regions": ["us-east-1","ap-south-1","us-east-2"]
}
```

APIs available 

| GET    /api/ec2/details           --> Gets all the EC2 instances along with state grouped by regions<br/> GET    /api/ec2/start              --> Start EC2<br/> GET    /api/ec2/stop              --> Stop EC2<br/> GET    /api/eip/details           --> Gets all unattached Elastic IPs<br/> GET    /api/eip/release          --> Release Elastic IP<br/> GET    /api/s3                        --> Gets all S3 buckets grouped by regions<br/> GET    /api/lb/details            --> Gets load balancers grouped by regions<br/> GET    /api/lb/delete             --> Delete Load Balancer<br/> GET    /api/dd/details           --> Gets DynamoDB tables |
| ------------------------------------------------------------ |

**References**

- https://www.reddit.com/r/golang/comments/n5ppx5/some_resources_that_have_helped_me_learn_golang/
- https://www.udemy.com/course/golang-for-devops-and-cloud-engineers/
- https://www.udemy.com/course/learn-how-to-code/