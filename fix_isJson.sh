#!/bin/bash
cat << 'INNER_EOF' > patch.diff
<<<<<<< SEARCH
export function isJson(str: string) {
    try {
        if (typeof JSON.parse(str) === 'object') {
            return true;
        }
    } catch {
        return false;
    }
}
=======
export function isJson(str: string): boolean {
    if (!str || typeof str !== 'string') {
        return false;
    }
    try {
        const obj = JSON.parse(str);
        if (obj && typeof obj === 'object') {
            return true;
        }
        return false;
    } catch {
        return false;
    }
}
>>>>>>> REPLACE
INNER_EOF
patch frontend/src/utils/misc.ts < patch.diff
