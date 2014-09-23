deb:
	GOOS=linux GOARCH=amd64 go build -o usr/bin/go-puppet-forge
	fpm -f -n go-puppet-forge -s dir -t deb \
		--workdir debian \
		--version `git describe --tags --long` \
		--deb-upstart debian/upstart/go-puppet-forge \
		--after-install debian/postinst usr/bin/
	rm -r usr

cent6:
	mkdir -p tmp
	GOOS=linux GOARCH=amd64 go build -o tmp/usr/local/bin/puppyforge
	mkdir -p tmp/etc/init.d/
	cp centos/6/puppyforge.conf tmp/etc/puppyforge.conf
	cp centos/6/puppyforge.init tmp/etc/init.d/puppyforge
	chmod +x tmp/etc/init.d/puppyforge
	#fpm -f -n puppyforge -s dir -t rpm \
	  --workdir tmp \
 		--version `git describe --tags --long`

