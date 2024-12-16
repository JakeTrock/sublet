"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SearchFeatures = void 0;
const fs_1 = __importDefault(require("fs"));
const lunr_1 = __importDefault(require("lunr"));
class SearchFeatures {
    constructor(channels, index) {
        this.channels = channels;
        this.entries = new Map();
        this.index = index;
    }
    search(value) {
        if (this.index && value.length >= 3) {
            // FIXME: do pagination here and not just harcoded limit of 30
            const results = this.index.search(value.split('.').join(' '))
                .splice(0, 30).map(result => result);
            const futures = results
                // translate to the cache key
                .map(r => this.transformPath(r.ref))
                // deduplicate
                .filter((item, i, ar) => ar.indexOf(item) === i)
                // and finally emit the download request
                .map(path => this.load(path));
            if (futures.length > 0) {
                forkJoin(...futures).subscribe(event => {
                    return results.map(r => {
                        const o = {
                            "path": r.ref,
                            "description": this.entries[r.ref].description,
                        };
                        return o;
                    });
                });
            }
            else {
                return results.map(r => `${r.ref} - ${this.entries[r.ref].description}`);
            }
        }
    }
    transformPath(value) {
        let parts = value.split('.');
        let beginning = parts.splice(0, 2);
        return beginning.join('_');
    }
    load() {
        this.channels.forEach(channel => {
            const index = JSON.parse(fs_1.default.readFileSync(`../indicies/options-${channel}.json`).toString());
            this.index = lunr_1.default.Index.load(index);
            Object.entries(index).forEach(entry => {
                this.entries.set(entry[0], entry[1]);
            });
        });
    }
}
exports.SearchFeatures = SearchFeatures;
