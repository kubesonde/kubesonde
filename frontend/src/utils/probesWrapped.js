/* eslint eqeqeq: 0, curly: 2 */
var parseAddress = exports.parseAddress = function (raw) {
    var port = null,
        address = null;
    if (raw[0] == '[') {
        port = raw.substring(raw.lastIndexOf(':') + 1);
        address = raw.substring(1, raw.indexOf(']'));
    } else if (raw.indexOf(':') != raw.lastIndexOf(':')) {
        port = raw.substring(raw.lastIndexOf(':') + 1);
        address = raw.substring(0, raw.lastIndexOf(':'));
    } else if (raw.indexOf(':') == raw.lastIndexOf(':') && raw.indexOf(':') != -1) {
        var parts = raw.split(':');
        port = parts[1];
        address = parts[0] || null;
    } else if (raw.indexOf('.') != raw.lastIndexOf('.')) {
        port = raw.substring(raw.lastIndexOf('.') + 1);
        address = raw.substring(0, raw.lastIndexOf('.'));
    }

    if (address && (address == '::' || address == '0.0.0.0')) {
        address = null;
    }

    return {
        port: port ? parseInt(port,10) : null,
        address: address
    };
};

const normalizeValues = function (item) {
    item.protocol = item.protocol.toLowerCase();
    var parts = item.local.split(':');
    item.local = parseAddress(item.local);
    item.remote = parseAddress(item.remote);

    if (item.protocol == 'tcp' && item.local.address && item.local.address.indexOf(':') !== -1) {
        item.protocol = 'tcp6';
    }

    if (item.pid == '-') {
        item.pid = 0;
    } else if (item.pid.indexOf('/') !== -1) {
        parts = item.pid.split('/');
        item.pid = parts.length > 1 ? parts[0] : 0;
    } else if (isNaN(item.pid)) {
        item.pid = 0;
    }

    item.pid = parseInt(item.pid,10);
    return item;
};

exports.linux = function (options) {
    options = options || {};
    var parseName = Boolean(options.parseName)

    return function (line, callback) {
        var parts = line.split(/\s/).filter(String);
        if (!parts.length || parts[0].match(/^(tcp|udp)/) === null) {
            return;
        }

        // NOTE: insert null for missing state column on UDP
        if (parts[0].indexOf('udp') === 0) {
            parts.splice(5, 0, null);
        }

        var name = '';
        var pid = parts.slice(6, parts.length).join(" ");
        if (parseName && pid.indexOf('/') > 0) {
            var pidParts = pid.split('/');
            pid = pidParts[0];
            name = pidParts.slice(1, pidParts.length).join('/');
        }

        var item = {
            protocol: parts[0],
            local: parts[3],
            remote: parts[4],
            state: parts[5],
            pid: pid
        };

        if (parseName) {
            item.processName = name;
        }

        return callback(normalizeValues(item));
    };
};
