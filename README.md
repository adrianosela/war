# war

A utility to run commands with assumed-role credentials in the environment.

### Usage: `war <role-arn-to-assume> <command> [args...]`

#### Examples

- `war arn:aws:iam::123456789012:role/MY-DEMO-ROLE aws ec2 describe-instances`
- `war arn:aws:iam::123456789012:role/MY-DEMO-ROLE aws iam list-roles`
- `war arn:aws:iam::123456789012:role/MY-DEMO-ROLE ${ANY_AWS_CLI_CMD}`
- `war arn:aws:iam::123456789012:role/MY-DEMO-ROLE ${ANY_NON_AWS_CLI_CMD_THAT_USES_AWS_CREDENTIALS}`

> Note: to test whether you are assuming the role, you can run aws sts get-caller-identity (which requires no IAM permissions to run)

```
$ war arn:aws:iam::123456789012:role/MY-DEMO-ROLE aws sts get-caller-identity
{
    "UserId": "AROHA7R44MA5HUOWH3YNXZ:WAR-1689374829999217000",
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/MY-DEMO-ROLE/WAR-1689374829999217000"
}
```
