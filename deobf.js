const fs = require('fs');
const path = require('path');
const parser = require('@babel/parser');
const traverse = require('@babel/traverse').default;
const t = require('@babel/types');
const generate = require('@babel/generator').default;
const vm = require('vm');

const infile = process.argv[2] || 'in.js';
const code = fs.readFileSync(path.resolve(infile), 'utf8');
const ast = parser.parse(code, {
	sourceType: 'unambiguous',
	plugins: ['jsx', 'classProperties', 'dynamicImport'],
});

let arrayName = null;
let strings = null;
let offset = null;
let lookupFunctionId = null;

traverse(ast, {
	VariableDeclarator(path) {
		const {id, init} = path.node;
		if (
			t.isIdentifier(id) &&
			t.isArrayExpression(init) &&
			init.elements.length > 0 &&
			init.elements.every(el => t.isStringLiteral(el))
		) {
			arrayName = id.name;
			strings = init.elements.map(el => el.value);
			path.stop();
		}
	},

});
traverse(ast, {
	FunctionDeclaration(path) {
		const {id, body} = path.node;
		if (!t.isIdentifier(id)) return;

		const stmts = body.body[0];

		if (stmts.type !== "ReturnStatement") return

		const expr = stmts.argument.expressions[0]

		if (expr.operator !== "=") return

		const num = expr.right.body.body[0].expression.right.right

		if (num.type !== "NumericLiteral") return;
		offset = num.value
		lookupFunctionId = id.name
		path.stop();
	}
})
// i couldnt bother 🙏
function a0_0x4494(_0x4b393f, _0x3cebfb) {
	return a0_0x4494 = function (_0x557bee, _0x449473) {
		_0x557bee = _0x557bee - 0x19d;
		var _0x4372bc = strings[_0x557bee];
		return _0x4372bc;
	}
		,
		a0_0x4494(_0x4b393f, _0x3cebfb);
}

(function (_0x3ba475, _0x540ee7) {
	var _0x34ba88 = a0_0x4494;
	while (!![]) {
		try {
			var _0x6573b6 = parseInt(_0x34ba88(0xa0a)) + -parseInt(_0x34ba88(0x442)) + parseInt(_0x34ba88(0x932)) * -parseInt(_0x34ba88(0x712)) + -parseInt(_0x34ba88(0x65e)) * parseInt(_0x34ba88(0x32b)) + parseInt(_0x34ba88(0x53e)) * -parseInt(_0x34ba88(0x801)) + parseInt(_0x34ba88(0x653)) + parseInt(_0x34ba88(0x3a1));
			if (_0x6573b6 === _0x540ee7)
				break;
			else
				_0x3ba475['push'](_0x3ba475['shift']());
		} catch (_0x424961) {
			_0x3ba475['push'](_0x3ba475['shift']());
		}
	}
}(strings, 0x3fef7))


console.log('Detected array:', arrayName);
console.log('Detected offset:', offset);
console.log('Detected lookup function:', lookupFunctionId);

if (!arrayName || !strings || offset === null || !lookupFunctionId) {
	console.log('Failed to detect necessary components');
	process.exit(1);
}


const aliases = new Set([lookupFunctionId]);

traverse(ast, {
	VariableDeclarator(path) {
		const {id, init} = path.node;
		if (t.isIdentifier(init) && aliases.has(init.name)) {
			if (t.isIdentifier(id)) {
				aliases.add(id.name);
			}
		}
	},
	AssignmentExpression(path) {
		const {left, right} = path.node;
		if (t.isIdentifier(right) && aliases.has(right.name)) {
			if (t.isIdentifier(left)) {
				aliases.add(left.name);
			}
		}
	},
});

traverse(ast, {
	CallExpression(path) {
		const {callee, arguments: args} = path.node;
		if (
			t.isIdentifier(callee) &&
			aliases.has(callee.name) &&
			args.length >= 1 &&
			t.isNumericLiteral(args[0])
		) {
			const num = args[0].value;
			const index = num - offset;
			if (index >= 0 && index < strings.length) {
				const str = strings[index];
				path.replaceWith(t.stringLiteral(str));
			} else {
				console.warn(`Index out of bounds: ${index} for call ${callee.name}(${num})`);
			}
		}
	},
});


const arrays = {};
const bindings = {};

traverse(ast, {
	VariableDeclarator(path) {
		const {id, init} = path.node;
		if (
			t.isIdentifier(id) &&
			t.isArrayExpression(init) &&
			init.elements.length > 0 &&
			init.elements.length < 2000 &&
			init.elements.every(el => t.isStringLiteral(el))

		) {
			const name = id.name;
			const binding = path.scope.getBinding(name);
			if (
				binding && binding.referencePaths.length === 2 && (
					(binding.referencePaths[0].parent.type == "MemberExpression"
						&& binding.referencePaths[1].parent.type == "CallExpression")
					|| (binding.referencePaths[1].parent.type == "MemberExpression"
						&& binding.referencePaths[0].parent.type == "CallExpression"))
			) {
				arrays[name] = init.elements.map(el => el.value);
				bindings[name] = binding;
			}
		}
	},
});

const arrayUsages = Object.fromEntries(
	Object.entries(bindings).map(([name, binding]) => {
		return [
			name,
			binding
				? binding.referencePaths.map(refPath => ({
					node: refPath.node,
					parent: refPath.parent
				}))
				: []
		];
	})
);

console.log('Detected arrays:', Object.keys(arrays));


function findMemberExpression(usages) {
	const identifierUsage = usages.find(u => u.parent.type === 'MemberExpression');
	if (!identifierUsage) return null;

	return usages.findIndex(u => u === identifierUsage);
}

function findFunctionExpression(usages) {
	const identifierUsage = usages.find(u => u.parent.type === 'CallExpression');
	if (!identifierUsage) return null;

	return usages.findIndex(u => u === identifierUsage);
}


for (const array of Object.keys(arrays)) {
	console.log(array)

	const usages = arrayUsages[array]
	const memberIndex = findMemberExpression(usages)
	const funcIndex = findFunctionExpression(usages)

	const offset = usages[memberIndex].parent.property.right.value
	const offsetFuncName = usages[funcIndex].parent.callee.body.body[1].init.declarations[0].init.name

	let code = `const ${array} = ${JSON.stringify(arrays[array])};\nfunction ${offsetFuncName}(arr){return ${array}[arr - ${offset}]};\n!`
	const funcNode = usages[funcIndex].parent.callee;

	funcNode.body.body = funcNode.body.body.filter(statement => {
		return !(
			statement.type === 'VariableDeclaration' &&
			statement.declarations.length === 1 &&
			statement.declarations[0].init &&
			statement.declarations[0].init.type === 'Identifier'
		);
	});
	code += generate(funcNode).code + `(${array})\n${array};`
	const context = {};
	const strings = vm.runInNewContext(code, context);

	const aliases = new Set([offsetFuncName]);

	traverse(ast, {
		VariableDeclarator(path) {
			const {id, init} = path.node;
			if (t.isIdentifier(init) && aliases.has(init.name)) {
				if (t.isIdentifier(id)) {
					aliases.add(id.name);
				}
			}
		},
		AssignmentExpression(path) {
			const {left, right} = path.node;
			if (t.isIdentifier(right) && aliases.has(right.name)) {
				if (t.isIdentifier(left)) {
					aliases.add(left.name);
				}
			}
		},
	});

	traverse(ast, {
		CallExpression(path) {
			const {callee, arguments: args} = path.node;
			if (
				t.isIdentifier(callee) &&
				aliases.has(callee.name) &&
				args.length >= 1 &&
				t.isNumericLiteral(args[0])
			) {
				const num = args[0].value;
				const index = num - offset;
				if (index >= 0 && index < strings.length) {
					const str = strings[index];
					path.replaceWith(t.stringLiteral(str));
				} else {
					console.warn(`Index out of bounds: ${index} for call ${callee.name}(${num})`);
				}
			}
		},
	});

}


const modifiedCode = generate(ast).code;
fs.writeFileSync('out.js', modifiedCode);
console.log('Modified code written to out.js');