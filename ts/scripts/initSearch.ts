import { promisify } from 'util';
import { exec as execCallback } from 'child_process';
import fs from 'fs';
import sqlite3 from 'sqlite3';

const exec = promisify(execCallback);

const channels = ["nixos-23.05", "nixos-23.11", "nixos-24.05", "nixos-unstable"];

async function get_options(channel: string) {
	const { stdout } = await exec(`nix-build '<nixpkgs/nixos/release.nix>' --no-out-link -A options -I nixpkgs=channel:${channel}`);
	const out_path = (await stdout).trim();

	return `${out_path}/share/doc/nixos/options.json`
}

function attributeByPath(obj: any, path: string[], def: any) {
	return path.reduce((acc: any, key: string) => ( acc[key] || def), obj)
}

function toSearchIndexObject(attrname: string, obj: any) {
	const i = {
		attribute_name: attrname,
		name: obj.name,
		license: attributeByPath(obj, ["meta", "license", "shortName"], "n/a"),
		long_description: obj.meta?.longDescription,
		example: attributeByPath(obj, ["example", "text"], "n/a"),
		default: attributeByPath(obj, ["default", "text"], "n/a"),
		type: attributeByPath(obj, ["type"], "n/a"),
	};

	if (obj.meta && obj.meta.longDescription)
		i.long_description = obj.meta.longDescription;

	if (obj.meta && obj.meta.license) {
		const l = obj.meta.license;
		if (l.shortName) {
			i.license = l.shortName;
		} else {
			i.license = l;
		}
	}
	return i;
}

channels.forEach(channel => {
	// Create SQLite full text search file for each channel
	get_options(channel).then(path => {
		const optionsData = JSON.parse(fs.readFileSync(path).toString());
		const optionsEntries = Object.entries(optionsData);

		const db = new sqlite3.Database(`options-${channel}.db`);
		db.serialize(() => {
			db.run('CREATE VIRTUAL TABLE options USING fts5(attribute_name, name, license, long_description, example, default, type)');

			const stmt = db.prepare('INSERT INTO options (attribute_name, name, license, long_description, example, default, type) VALUES (?, ?, ?, ?, ?, ?, ?)');
			for (const [attrName, optionData] of optionsEntries) {
				const searchEntry = toSearchIndexObject(attrName, optionData);
				stmt.run(
					searchEntry.attribute_name,
					searchEntry.name,
					searchEntry.license,
					searchEntry.long_description,
					searchEntry.example,
					searchEntry.default,
					searchEntry.type
				);
			}
			stmt.finalize();
		});
		db.close();
	});
});
