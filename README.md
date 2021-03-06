# bristol.xyz
The codebase for bristol.xyz.

## Requirements
The requirements for hosting this are the following:
- A MongoDB database. This stores all of the records which are used on the website.
- A S3-compatible storage solution. This stores all the files which are used on the website. In production we use DigitalOcean Spaces. For development, you can use a solution such as Amazon S3 on AWS free tier if you do not already have a Spaces subscription.

## Configuration

Configuration can be quite simply done in a few steps:

### Environment Variables
To configure bristol.xyz, you will need to set the following environment variables:

#### Initialisation
- `INITIAL_USER` - A user in the format of `email:password` which will be initially created with administrator permissions if this is a initial install (there are no items in the `users` collection).

#### Sentry Configuration
- `SENTRY_DSN` - The DSN which is in use for Sentry.

#### MongoDB Configuration
- `MONGODB_URI` - The connection URI for your MongoDB instance. Defaults to `mongodb://localhost:27017`.
- `MONGODB_DATABASE` - The database to use for MongoDB. Defaults to `bristolxyz`.

#### S3 Configuration
- `S3_ENDPOINT` - The endpoint for the S3 provider such as Amazon S3 or DigitalOcean Spaces.
- `S3_REGION` - The S3 region.
- `S3_BUCKET` - The S3 bucket name.
- `CDN_HOSTNAME` - The hostname for the CDN which is running from the bucket. For example, `cdn.bristol.xyz`.
- `AWS_SECRET_ACCESS_KEY` - The AWS secret access key for the bucket which we are using for the CDN.
- `AWS_ACCESS_KEY_ID` - The AWS access key ID for the bucket.

#### Mailgun Configuration
- `MAILGUN_DOMAIN` - The domain used for Mailgun.
- `MAILGUN_KEY` - The private key for your Mailgun instance.
- `FROM_ADDRESS` - The e-mail which e-mails are from (in the format `Name <email>`).

### Setting up your local environment/production deployment
Now we have the environment variables listed down, we can continue the setup:

#### Local Development
For working on the project locally, you can simply make a `.env` file with the contents above in. This is ignored by Git in this project, so you can rest assure that you will not accidentally commit this. Ensure you have Go 1.13+ installed and then you can run `go run .`.

#### Production deployment
In production, we use the Dockerfile which is in the root of this project. From here, you can pass all of the environment variables in from above, making it a pretty painless deploy.
