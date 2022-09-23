There are two components to the system:

    - an authenticator
    - a discriminator

Ideally, they are deployed at different networks (to avoid correlation).

Presenting the right credentials to the listener opens a window in the
discriminator. Credentials for authentication are cryptographically signed and
are assumed to be distributed off-band.

Discriminator should serve a legitimate resource as a default.

The discriminator queries the authenticator service to switch traffic for requests.

- Every server generates a key pair.
- Each server can generate a bunch of one-time tokens, signed by its private
  key.
- The client sends credentials to the authenticator.
- Credentials can be obfuscated; authenticator will reassemble them according
  to some schema. There's room for improvement in this step.
- Upon successful authentication, the authenticator service will return a valid
  endpoint, and a one-time token valid for the discriminator.
- Additional rules can be imposed on the access to the discriminator (check for
  timing, trust score of the requester, etc).


