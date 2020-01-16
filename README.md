# MCCS Alpha

Looking to try it out? See the [getting started](#how-to-start) instructions.

See our [roadmap](https://github.com/ic3network/mccs#roadmap) of what we are working on next.

## Overview

What this alpha software is:

- A prototype written specifically for the [Open Credit Network](https://opencredit.network) (OCN)
- A "throw away", proof of concept to allow OCN to experiment with running a mutual credit system

We are making this code public to show our commitment to free and open source software, and to signal our intention to develop mutual credit software that will be freely available to anyone who wishes to implement a mutual credit trading system.

This alpha version of the software will not be maintained over the long-term because it was built as a proof of concept specifically for OCN. Although you are free to copy this software and modify it for your own use, you might want to wait for the next release of MCCS which will expose an API to all of the functionality currently available in the alpha version of the MCCS web application. By providing an API, developers who want to create their own version of MCCS will have complete flexibility to implement it in whatever way they choose. This means developers can present MCCS in any language, optimize it for whatever devices their user base prefer, develop a mobile app, etc.

## Main Functions

There are four main functions that the MCCS web application provides:

1. **Manage accounts** - create and modify user accounts and related business details
2. **Find businesses** - view and search for businesses based on what they sell and need
3. **Transfer mutual credits** - create and complete/cancel mutual credit (MC) transfers between businesses
4. **Review transfer activity** - view pending and completed MC transfers

Users access MCCS through either a desktop or mobile web browser. 

### Manage Accounts

Individual users can create an account in MCCS by providing an email address, creating a password and adding some other details about themselves and their business (a "business" need not be a formal business; it could simply be a list of their skills that they are willing to offer to other participants in the network). Account management includes the expected features like:

- Create account
- Login/logout
- View and modify one's user & business details
- Reset a lost password
- Change password

### Find Businesses

MCCS enables businesses to list their "offers" (goods and services it can provide to the network) and "wants" (goods and services it would like to procure from the network), and to see the offers and wants of the other participants. There is a matching algorithm that looks for corresponding wants that match a business' offers and offers that match a business' wants.

From the users' perspective, once they have listed their offers and wants, MCCS will propose potential trades to them with other users. These potential trades are displayed to them on the dashboard immediately after logging in, and users can opt into a daily email notification that displays all of the matches to their offers and wants that occured in the previous 24 hours.

MCCS' matching logic searches for exact matches, wildcard matches and fuzzy matches, so a search for "vegetables" will not only turn up "vegetables" (exact match) but also "organic-vegetables" (wildcard match) and "vegetalbes" (fuzzy match, which can pick up typos or spelling variations).

Aside from the tagging feature on offers and wants, MCCS also offers a standard directory (like the Yellow Pages). An admin running MCCS can assign various categories to any business and the users can browse these categories alphabetically, just like in a telephone directory.

Businesses are able to send contact emails to each other to start a dialog about a trade. The business receiving the first email does not have to expose its email address to the party contacting them. Only if the business replies will its email be exposed as the conversation continues between them through regular email.

Finally, users can flag specific businesses as "favorites" so they can be easily found.

### Transfer Mutual Credits

Each business in MCCS will have a transaction account that it can pay from or accept payment to in order to facilitate transfers between MCCS users. The sum of those increases and decreases will yield the current balance of that business' account. IOW, it is the same functionality you would expect from a [basic bank accounting system](https://medium.com/@RobertKhou/double-entry-accounting-in-a-relational-database-2b7838a5d7f8).

The following rules will be enforced by MCCS:

- Every increase in one account must be offset by an equal decrease in another account (or set of accounts in the case of collecting fees). So if Account A is decreased by 50, Account B will be increased by 50 (or Account B could be increased by 49 and Account C by 1).
- Given the above, the sum total of all transfers, and the overall balance in the system must always equal zero. These are the two foundational principles of double-entry bookkeeping and they must be enforced at all times.
- In order to start the transfer process, some accounts will need to have a negative balance (see Account A above), but not all accounts should be allowed to have a negative balance. MCCS admins will need to make a decision about which businesses may have an account whose balance can be negative.
- There needs to be a limit to the positive and negative balances for accounts (maxPositiveBalance and maxNegativeBalance). A transfer can only be processed if the resulting balances of all accounts involved in the transfer do not exceed these limits. Continuing the example above, if the maxNegativeBalance for Account A is set to -40, the transaction above would not be processed because the balance of Account A would be exceeding this limit by going to -50. And if the maxNegativeBalance for an account is set to 0, then that account is not allowed to have a negative balance at all.

#### Initiating Debits or Credits

Users will be able to initiate transfers that result in either a **debit from** or a **credit to** their account. This is unusual given that most online payment systems only allow the former. A practical application of this ability to initiate a credit to one's account is that a supplier could setup the payment from their customer within MCCS. The customer would login, see the invoice link in the payment description and if they agrees, confirm the debit of funds from their account. 

#### Two-Step Transfers

One of the objectives behind a mutual credit system is to enable participants to vet transfers before the funds are added to or, in the case of payee-initiated transfers, removed from their accounts. This requires that a transfer must first be authorized by by both the sender and receiver of the credits before the transfer is completed and the accounts are debited and credited accordingly.

To enable this two step transfer process, state has been introduced transfers, so a transfer can be in one of only three possible states at any given moment:

1. **Initiated** - a transfer has been initiated by either the payer or payee of the credits.
2. **Completed** - a transfer is completed because the payee/payer have authorized the credit to/debit from their account that was proposed in the initiated transfer.
3. **Cancelled** - a transfer can be cancelled because (a) the initiator cancels it before the receiver has a chance to accept it, (b) the receiver can decide to reject the transfer, or (c) the transfer is not allowed by the system due to balance limits or other reasons even if the receiver accepts it.

Only transfers in the completed state will affect the balances of the payer's and payee's accounts, and therefore be shown in the transaction history of their accounts. Initiated payments will be shown in the Dashboard summary of pending transfers for both the payer and payee, and cancelled payments will trigger an email to both the payer and payee explaining why the payment was cancelled.

### Review Transfer Activity

Completed transfers are presented as a statement like you would see in an online banking application, and include the current credit balance of the account.

Initiated transfers are presented to the users. If users created a transfer, they have the option to cancel it before the receiver of the transfer (which could be a credit to or debit from their account) approves or rejects the transfer. Likewise if users receive a transfer, they can decide whether to approve or reject it. Approved transfers are set to the completed state, rejected (by the receiver) and cancelled (by the initiator) transfers are set to the cancelled state.

## Demo Site

A demo site of MCCS is available at:  
https://alpha.ic3.dev

You can create a new account by signing up:  
https://alpha.ic3.dev/signup

Or just login and use one of the demo accounts:  
user =jdoe1@dev.null (or jdoe2, jdoe3 ... up to jdoe5)  
pass = password (all lower case)

And if you want to see the admin side, here's the login:   
https://alpha.ic3.dev/admin/login  
user = admin1@dev.null or admin2@dev.null  
pass = password

There are 200 businesses made up of dummy data in the database already.

The demo site will be reset from time to time, so any test accounts you create will eventually be wiped.


## How to Start

Basic requirements: Go version 1.13+, Docker and Docker Compose ([see all requirements](#requirements))

1. Use the [example file](configs/seed.yaml) to create `configs/development.yaml` and change the following parameters
    ```
    env: development
    
    psql:
      host: postgres

    mongo:
        url: mongodb://mongo:27017

    es:
        url: http://es01:9200
    ```
1. Generate JSON Web Token public and private keys
    1. Generate private key
        ```
        openssl genrsa -out private.pem 2048
        ```
    1. Extract public key from it
        ```
        openssl rsa -in private.pem -pubout > public.pem
        ```
    1. Copy/paste the private and public keys into `configs/development.yaml`
1. Run the app
    ```
    make run
    ```
1. Seed the DB with some test data
    ```
    make seed
    ```
1. Visit the website
    ```
    http://localhost:8080/signup
    http://localhost:8080/login
    ```
1. Login as admin (in a private browsing tab or another browser)
    ```
    http://localhost:8080/admin/login
    user = admin1@dev.null
    pass = password
    ```
1. View/query Elasticsearch data with Kibana
    ```
    http://localhost:5601
    ```

## Requirements

**Software**

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

The MCCS web app is written in Go, and it uses Docker Compose to orchestrate itself and its dependencies as Docker containers.

**App Dependencies**

- [MongoDB](https://en.wikipedia.org/wiki/MongoDB) - the database used to store user and business directory data
- [PostgreSQL](https://www.postgresql.org/) - the database used to store mutual credit transfer data
- [Elasticsearch](https://en.wikipedia.org/wiki/Elasticsearch) - the search engine for searching user/business data

These dependencies are installed automatically by Docker Compose.

**External Dependencies**

- [reCAPTCHA](https://www.google.com/recaptcha) - protects against automated hacking attempts when running in production mode
- [Sendgrid](https://sendgrid.com/) - email delivery provider for sending emails generated by MCCS (e.g., welcome and password reset emails)

Both reCAPTCHA and Sendgrid do not need to be setup to run the app in development mode.
