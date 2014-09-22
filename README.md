# puppyforge

This is a minimalistic Puppet Forge implementation that implements the v3 forge API based on modules stored on disk. No database is required, just read access to the module files.

### Installation
Copy your module .tar.gz files to a directory. `/var/lib/puppyforge/modules/` for example:

    /var/lib/puppyforge/modules/puppetlabs-apache-1.1.0.tar.gz

### Running puppyforge
#### Via the command-line
Running puppyforge from the command line is useful for development and testing.

    $ ./puppyforge --port 8080 --modulepath /var/lib/puppyforge/modules

#### Via a service
Usually you will want to run puppyforge via a service. Building the packages for your operating system provides the correct service.

    $ cat /etc/puppyforge.conf
    $ service puppyforge start
        
#### Usage with Puppet
Once your puppyforge is running, you can configure puppet to use it via the `module_repository` option in the puppet config file.

    module_repository=http://my-forge.com/

Or directly on command line

    $ puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath modules
