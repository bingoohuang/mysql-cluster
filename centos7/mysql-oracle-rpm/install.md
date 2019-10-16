1. [Installing MySQL on Linux Using the MySQL Yum Repository](https://dev.mysql.com/doc/mysql-repo-excerpt/5.7/en/linux-installation-yum-repo.html)
1. Goo to the [Download MySQL Yum Repository page](https://dev.mysql.com/downloads/repo/yum/) in the MySQL Developer Zone.
1. Select and download the release package for your platform.
1. `sudo yum localinstall -y  mysql80-community-release-el7-3.noarch.rpm`
    ```bash
    [vagrant@10 haproxy-1.8.21]$ sudo yum localinstall -y  mysql80-community-release-el7-3.noarch.rpm
    Loaded plugins: fastestmirror
    Examining mysql80-community-release-el7-3.noarch.rpm: mysql80-community-release-el7-3.noarch
    Marking mysql80-community-release-el7-3.noarch.rpm to be installed
    Resolving Dependencies
    --> Running transaction check
    ---> Package mysql80-community-release.noarch 0:el7-3 will be installed
    --> Finished Dependency Resolution

    Dependencies Resolved

    ======================================================================================================================
    Package                           Arch           Version       Repository                                       Size
    ======================================================================================================================
    Installing:
    mysql80-community-release         noarch         el7-3         /mysql80-community-release-el7-3.noarch          31 k

    Transaction Summary
    ======================================================================================================================
    Install  1 Package

    Total size: 31 k
    Installed size: 31 k
    Downloading packages:
    Running transaction check
    Running transaction test
    Transaction test succeeded
    Running transaction
    Installing : mysql80-community-release-el7-3.noarch                                                             1/1
    Verifying  : mysql80-community-release-el7-3.noarch                                                             1/1

    Installed:
    mysql80-community-release.noarch 0:el7-3

    Complete!
    ```
1. check that the MySQL Yum repository has been successfully added `yum repolist enabled | grep "mysql.*-community.*"`
    ```bash
    [vagrant@10 haproxy-1.8.21]$ yum repolist enabled | grep "mysql.*-community.*"
    mysql-connectors-community/x86_64       MySQL Connectors Community          128
    mysql-tools-community/x86_64            MySQL Tools Community               100
    mysql80-community/x86_64                MySQL 8.0 Community Server          145
    ```
1. Selecting a Release Series
    ```bash
    [vagrant@10 haproxy-1.8.21]$  yum repolist all | grep mysql
    mysql-cluster-7.5-community/x86_64 MySQL Cluster 7.5 Community    disabled
    mysql-cluster-7.5-community-source MySQL Cluster 7.5 Community -  disabled
    mysql-cluster-7.6-community/x86_64 MySQL Cluster 7.6 Community    disabled
    mysql-cluster-7.6-community-source MySQL Cluster 7.6 Community -  disabled
    mysql-cluster-8.0-community/x86_64 MySQL Cluster 8.0 Community    disabled
    mysql-cluster-8.0-community-source MySQL Cluster 8.0 Community -  disabled
    mysql-connectors-community/x86_64  MySQL Connectors Community     enabled:   128
    mysql-connectors-community-source  MySQL Connectors Community - S disabled
    mysql-tools-community/x86_64       MySQL Tools Community          enabled:   100
    mysql-tools-community-source       MySQL Tools Community - Source disabled
    mysql-tools-preview/x86_64         MySQL Tools Preview            disabled
    mysql-tools-preview-source         MySQL Tools Preview - Source   disabled
    mysql55-community/x86_64           MySQL 5.5 Community Server     disabled
    mysql55-community-source           MySQL 5.5 Community Server - S disabled
    mysql56-community/x86_64           MySQL 5.6 Community Server     disabled
    mysql56-community-source           MySQL 5.6 Community Server - S disabled
    mysql57-community/x86_64           MySQL 5.7 Community Server     disabled
    mysql57-community-source           MySQL 5.7 Community Server - S disabled
    mysql80-community/x86_64           MySQL 8.0 Community Server     enabled:   145
    mysql80-community-source           MySQL 8.0 Community Server - S disabled
    [vagrant@10 haproxy-1.8.21]$ sudo yum-config-manager --disable mysql80-community
    [vagrant@10 haproxy-1.8.21]$ sudo yum-config-manager --enable mysql57-community
    Loaded plugins: fastestmirror
    ========================================================================================================== repo: mysql57-community ===========================================================================================================
    [mysql57-community]
    async = True
    bandwidth = 0
    base_persistdir = /var/lib/yum/repos/x86_64/7
    baseurl = http://repo.mysql.com/yum/mysql-5.7-community/el/7/x86_64/
    cache = 0
    cachedir = /var/cache/yum/x86_64/7/mysql57-community
    check_config_file_age = True
    compare_providers_priority = 80
    cost = 1000
    deltarpm_metadata_percentage = 100
    deltarpm_percentage =
    enabled = 1
    enablegroups = True
    exclude =
    failovermethod = priority
    ftp_disable_epsv = False
    gpgcadir = /var/lib/yum/repos/x86_64/7/mysql57-community/gpgcadir
    gpgcakey =
    gpgcheck = True
    gpgdir = /var/lib/yum/repos/x86_64/7/mysql57-community/gpgdir
    gpgkey = file:///etc/pki/rpm-gpg/RPM-GPG-KEY-mysql
    hdrdir = /var/cache/yum/x86_64/7/mysql57-community/headers
    http_caching = all
    includepkgs =
    ip_resolve =
    keepalive = True
    keepcache = False
    mddownloadpolicy = sqlite
    mdpolicy = group:small
    mediaid =
    metadata_expire = 21600
    metadata_expire_filter = read-only:present
    metalink =
    minrate = 0
    mirrorlist =
    mirrorlist_expire = 86400
    name = MySQL 5.7 Community Server
    old_base_cache_dir =
    password =
    persistdir = /var/lib/yum/repos/x86_64/7/mysql57-community
    pkgdir = /var/cache/yum/x86_64/7/mysql57-community/packages
    proxy = False
    proxy_dict =
    proxy_password =
    proxy_username =
    repo_gpgcheck = False
    retries = 10
    skip_if_unavailable = False
    ssl_check_cert_permissions = True
    sslcacert =
    sslclientcert =
    sslclientkey =
    sslverify = True
    throttle = 0
    timeout = 30.0
    ui_id = mysql57-community/x86_64
    ui_repoid_vars = releasever,
    basearch
    username =

    [vagrant@10 haproxy-1.8.21]$  yum repolist all | grep mysql
    Failed to set locale, defaulting to C
    mysql-cluster-7.5-community/x86_64 MySQL Cluster 7.5 Community    disabled
    mysql-cluster-7.5-community-source MySQL Cluster 7.5 Community -  disabled
    mysql-cluster-7.6-community/x86_64 MySQL Cluster 7.6 Community    disabled
    mysql-cluster-7.6-community-source MySQL Cluster 7.6 Community -  disabled
    mysql-cluster-8.0-community/x86_64 MySQL Cluster 8.0 Community    disabled
    mysql-cluster-8.0-community-source MySQL Cluster 8.0 Community -  disabled
    mysql-connectors-community/x86_64  MySQL Connectors Community     enabled:   128
    mysql-connectors-community-source  MySQL Connectors Community - S disabled
    mysql-tools-community/x86_64       MySQL Tools Community          enabled:   100
    mysql-tools-community-source       MySQL Tools Community - Source disabled
    mysql-tools-preview/x86_64         MySQL Tools Preview            disabled
    mysql-tools-preview-source         MySQL Tools Preview - Source   disabled
    mysql55-community/x86_64           MySQL 5.5 Community Server     disabled
    mysql55-community-source           MySQL 5.5 Community Server - S disabled
    mysql56-community/x86_64           MySQL 5.6 Community Server     disabled
    mysql56-community-source           MySQL 5.6 Community Server - S disabled
    mysql57-community/x86_64           MySQL 5.7 Community Server     enabled:   384
    mysql57-community-source           MySQL 5.7 Community Server - S disabled
    mysql80-community/x86_64           MySQL 8.0 Community Server     disabled
    mysql80-community-source           MySQL 8.0 Community Server - S disabled
    ```
1. Besides using yum-config-manager or the dnf config-manager command, you can also select a release series by editing manually the /etc/yum.repos.d/mysql-community.repo file. This is a typical entry for a release series' subrepository in the file:
    ```
    [mysql57-community]
    name=MySQL 5.7 Community Server
    baseurl=http://repo.mysql.com/yum/mysql-5.7-community/el/6/$basearch/
    enabled=1
    gpgcheck=1
    gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-mysql
    ```
1. 下载方式安装  `sudo yum install -y --downloadonly --downloaddir=. mysql-community-server`
    ```bash
    [vagrant@10 mysqlrpm]$ ll
    total 265804
    -rw-r--r--. 1 root root     24744 Nov 25  2015 libaio-0.3.109-13.el7.x86_64.rpm
    -rw-r--r--. 1 root root  45109364 Oct 10 16:06 mysql-community-client-5.7.28-1.el7.x86_64.rpm
    -rw-r--r--. 1 root root    318768 Oct 10 16:06 mysql-community-common-5.7.28-1.el7.x86_64.rpm
    -rw-r--r--. 1 root root   4374364 Oct 10 16:07 mysql-community-libs-5.7.28-1.el7.x86_64.rpm
    -rw-r--r--. 1 root root   1353312 Oct 10 16:07 mysql-community-libs-compat-5.7.28-1.el7.x86_64.rpm
    -rw-r--r--. 1 root root 208694824 Oct 10 16:07 mysql-community-server-5.7.28-1.el7.x86_64.rpm
    -rw-r--r--. 1 root root    312968 Aug 22 21:36 net-tools-2.0-0.25.20131004git.el7.x86_64.rpm
    -rw-r--r--. 1 root root   8358460 Jan 24  2019 perl-5.16.3-294.el7_6.x86_64.rpm
    -rw-r--r--. 1 root root     19672 Jul  4  2014 perl-Carp-1.26-244.el7.noarch.rpm
    -rw-r--r--. 1 root root   1545440 Jul  4  2014 perl-Encode-2.51-7.el7.x86_64.rpm
    -rw-r--r--. 1 root root     29092 Jul  4  2014 perl-Exporter-5.68-3.el7.noarch.rpm
    -rw-r--r--. 1 root root     27088 Jul  4  2014 perl-File-Path-2.09-2.el7.noarch.rpm
    -rw-r--r--. 1 root root     57680 Jul  4  2014 perl-File-Temp-0.23.01-3.el7.noarch.rpm
    -rw-r--r--. 1 root root     78236 Jul  4  2014 perl-Filter-1.49-3.el7.x86_64.rpm
    -rw-r--r--. 1 root root     57176 Apr 25  2018 perl-Getopt-Long-2.40-3.el7.noarch.rpm
    -rw-r--r--. 1 root root     39292 Jul  4  2014 perl-HTTP-Tiny-0.033-3.el7.noarch.rpm
    -rw-r--r--. 1 root root     84468 Jul  4  2014 perl-PathTools-3.40-5.el7.x86_64.rpm
    -rw-r--r--. 1 root root     52516 Jan 24  2019 perl-Pod-Escapes-1.04-294.el7_6.noarch.rpm
    -rw-r--r--. 1 root root     88756 Jul  4  2014 perl-Pod-Perldoc-3.20-4.el7.noarch.rpm
    -rw-r--r--. 1 root root    221216 Jul  4  2014 perl-Pod-Simple-3.28-4.el7.noarch.rpm
    -rw-r--r--. 1 root root     27436 Jul  4  2014 perl-Pod-Usage-1.63-3.el7.noarch.rpm
    -rw-r--r--. 1 root root     36808 Jul  4  2014 perl-Scalar-List-Utils-1.27-248.el7.x86_64.rpm
    -rw-r--r--. 1 root root     49812 Nov 20  2016 perl-Socket-2.010-4.el7.x86_64.rpm
    -rw-r--r--. 1 root root     78888 Jul  4  2014 perl-Storable-2.45-3.el7.x86_64.rpm
    -rw-r--r--. 1 root root     14056 Jul  4  2014 perl-Text-ParseWords-3.29-4.el7.noarch.rpm
    -rw-r--r--. 1 root root     46304 Jul  4  2014 perl-Time-HiRes-1.9725-3.el7.x86_64.rpm
    -rw-r--r--. 1 root root     24792 Jul  4  2014 perl-Time-Local-1.2300-2.el7.noarch.rpm
    -rw-r--r--. 1 root root     19244 Jul  4  2014 perl-constant-1.27-2.el7.noarch.rpm
    -rw-r--r--. 1 root root    704872 Jan 24  2019 perl-libs-5.16.3-294.el7_6.x86_64.rpm
    -rw-r--r--. 1 root root     44780 Jan 24  2019 perl-macros-5.16.3-294.el7_6.x86_64.rpm
    -rw-r--r--. 1 root root     12592 Jul  4  2014 perl-parent-0.225-244.el7.noarch.rpm
    -rw-r--r--. 1 root root    114320 Jul  4  2014 perl-podlators-2.5.1-3.el7.noarch.rpm
    -rw-r--r--. 1 root root     50392 Jul  4  2014 perl-threads-1.87-4.el7.x86_64.rpm
    -rw-r--r--. 1 root root     39868 Jul  4  2014 perl-threads-shared-1.43-6.el7.x86_64.rpm
    ```
