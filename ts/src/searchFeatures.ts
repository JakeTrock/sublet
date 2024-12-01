import fs from 'fs';
import lunr from 'lunr';


export class SearchFeatures {
  index: lunr.Index;

  entries: Map<string, any>;


  constructor(private channels: string[], index: lunr.Index) {
    this.entries = new Map();
    this.index = index;
  }

  public search(value: string): void {
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
        })
      } else {
        return results.map(r => `${r.ref} - ${this.entries[r.ref].description}`);
      }

    }
  }

  transformPath(value: string): string {
    let parts = value.split('.');
    let beginning = parts.splice(0, 2);
    return beginning.join('_');
  }


  load() {
    this.channels.forEach(channel => {
      const index = JSON.parse(fs.readFileSync(`../indicies/options-${channel}.json`).toString());
      this.index = lunr.Index.load(index);
      Object.entries(index).forEach(entry => {
        this.entries.set(entry[0], entry[1]);
      });
    });
  }
}