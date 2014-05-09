Name: hydra
Version: 3
Release: 0
Summary: Hydra
Source0: hydra.tar.gz
License: MIT
Group: custom
URL: https://github.com/innotech/hydra
BuildArch: noarch
BuildRoot: %{_tmppath}/%{name}-buildroot
# Requires: python-psutil
%description
Hydra is multi-cloud application discovery, management and balancing service.
Hydra attempts to ease the routing and balancing burden from servers and delegate it on the client (browser, mobile app, etc).
%prep
%setup -q
%build
%install
install -m 0755 -d $RPM_BUILD_ROOT/usr/local/hydra

install -m 0755 -d $RPM_BUILD_ROOT/etc/init.d
install -m 0755 hydra-init.d.sh $RPM_BUILD_ROOT/etc/init.d/hydra

install -m 0755 -d $RPM_BUILD_ROOT/etc/hydra
install -m 0644 hydra.conf $RPM_BUILD_ROOT/etc/hydra/hydra.conf
install -m 0644 apps-example.json $RPM_BUILD_ROOT/etc/hydra/apps-example.json
%clean
rm -rf $RPM_BUILD_ROOT
%post
echo   You should edit config file /etc/hydra/hydra.conf
echo   When finished, you may want to run \"update-rc.d hydra defaults\"
%files
%dir /etc/hydra
/usr/local/hydra/hydra
/etc/hydra/hydra.conf
/etc/hydra/apps-example.json
/etc/init.d/hydra
