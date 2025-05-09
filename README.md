# Loan Service
Loan Service provides API to manage Loan lifecyle from create, approve, invest until disburse.
It also provide API to retrieve the loans.

## Setup
To run Loan Service, the following infrastructure is required
1. Database mysql, the script to create tables is available in folder `migration`
2. SMTP server, use smtp server like `gmail` (not mandatory since sending email is put in asynchronous process)

The configuration for the infrastructure and http server are put in file `config.json`
