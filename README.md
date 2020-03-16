# bristol.xyz
The codebase for bristol.xyz.

## Requirements
The requirements for hosting this are the following:
- A Redis database. This stores all of the records which are used on the website.
- A S3-compatible storage solution. This stores all the files which are used on the website. In production we use DigitalOcean Spaces. For development, you can use a solution such as Amazon S3 on AWS free tier if you do not already have a Spaces subscription.

## Configuration

Configuration can be quite simply done in a few steps:

### Configuring the bucket
When you configure your Amazon S3/DigitalOcean Spaces bucket, you will want sure that the public can read everything in the bucket, but not list everything in the bucket. Instructions on how to do this can vary between providers. From here, you will want to copy everything inside the `bucket_base` folder in this repository and put it inside the bucket. Your bucket is ready!

### Environment Variables
To configure bristol.xyz, you will need to set the following environment variables:

#### Initialisation
- `INITIAL_USER` - A user in the format of `email:password` which will be initially created with administrator permissions if this is a initial install (there are no keys in the Redis database).

#### Sentry Configuration
- `SENTRY_DSN` - The DSN which is in use for Sentry. A blank DSN will result in errors being logged to the console.

#### Redis Configuration
- `REDIS_HOST` - The Redis host. Defaults to `localhost:6379`.
- `REDIS_PASSWORD` - The password for the Redis database. Defaults to none.

#### S3 Configuration
- `S3_ENDPOINT` - The endpoint for the S3 provider such as Amazon S3 or DigitalOcean Spaces.
- `S3_BUCKET` - The S3 bucket name.
- `CDN_HOSTNAME` - The hostname for the CDN which is running from the bucket. For example, `cdn.bristol.xyz`.
- `AWS_SECRET_ACCESS_KEY` - The AWS secret access key for the bucket which we are using for the CDN.
- `AWS_ACCESS_KEY_ID` - The AWS access key ID for the bucket.

### Setting up your local environment/production deployment
Now we have the environment variables listed down, we can continue the setup:

#### Local Development
For working on the project locally, you can simply make a `.env` file with the contents above in. This is ignored by Git in this project, so you can rest assure that you will not accidentally commit this. Ensure you have a version of node, nvm and Go 1.14+ installed and then you can do one of the following:
- **Run `./build.sh`** - This will build a binary which you can run in the folder of this repository. Note that you will lose all debugging abilities doing this.
- **Use the built in Visual Studio Code configuration** - This is probably the best option if you do not mind using VS Code.
- **Manually configure your IDE** - The basic steps that your IDE will want to do is run `cd frontend && npm i && npm run build` before debugging the Go project. Doing this can vary between IDE's.

#### Production deployment
In production, we use the Dockerfile which is in the root of this project. From here, you can pass all of the environment variables in from above, making it a pretty painless deploy.
