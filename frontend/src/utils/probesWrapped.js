/* eslint eqeqeq: 0, curly: 2 */

export function parseAddress(raw) {
    let port = null;
    let address = null;
  
    if (raw[0] === '[') {
      port = raw.substring(raw.lastIndexOf(':') + 1);
      address = raw.substring(1, raw.indexOf(']'));
    } else if (raw.indexOf(':') !== raw.lastIndexOf(':')) {
      port = raw.substring(raw.lastIndexOf(':') + 1);
      address = raw.substring(0, raw.lastIndexOf(':'));
    } else if (raw.indexOf(':') === raw.lastIndexOf(':') && raw.indexOf(':') !== -1) {
      const parts = raw.split(':');
      port = parts[1];
      address = parts[0] || null;
    } else if (raw.indexOf('.') !== raw.lastIndexOf('.')) {
      port = raw.substring(raw.lastIndexOf('.') + 1);
      address = raw.substring(0, raw.lastIndexOf('.'));
    }
  
    if (address && (address === '::' || address === '0.0.0.0')) {
      address = null;
    }
  
    return {
      port: port ? parseInt(port, 10) : null,
      address
    };
  }
  
  function normalizeValues(item) {
    item.protocol = item.protocol.toLowerCase();
    let parts = item.local.split(':');
    item.local = parseAddress(item.local);
    item.remote = parseAddress(item.remote);
  
    if (item.protocol === 'tcp' && item.local.address && item.local.address.includes(':')) {
      item.protocol = 'tcp6';
    }
  
    if (item.pid === '-') {
      item.pid = 0;
    } else if (item.pid.includes('/')) {
      parts = item.pid.split('/');
      item.pid = parts.length > 1 ? parts[0] : 0;
    } else if (isNaN(item.pid)) {
      item.pid = 0;
    }
  
    item.pid = parseInt(item.pid, 10);
    return item;
  }
  
  export function linux(options = {}) {
    const parseName = Boolean(options.parseName);
  
    return function (line, callback) {
      const parts = line.split(/\s/).filter(Boolean);
      if (!parts.length || parts[0].match(/^(tcp|udp)/) === null) {
        return;
      }
  
      // NOTE: insert null for missing state column on UDP
      if (parts[0].startsWith('udp')) {
        parts.splice(5, 0, null);
      }
  
      let name = '';
      let pid = parts.slice(6).join(' ');
      if (parseName && pid.includes('/')) {
        const pidParts = pid.split('/');
        pid = pidParts[0];
        name = pidParts.slice(1).join('/');
      }
  
      const item = {
        protocol: parts[0],
        local: parts[3],
        remote: parts[4],
        state: parts[5],
        pid
      };
  
      if (parseName) {
        item.processName = name;
      }
  
      return callback(normalizeValues(item));
    };
  }
  