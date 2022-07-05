
const execSync = require('child_process').execSync;


const submodule_sync = [
    ['view/extensions/analytics', '1.0.0'],
    ['view/extensions/auth-github', 'v0.1.0-beta.13'],
    ['view/extensions/diff', '1.1.1'],
    ['view/extensions/flarum-pipetables', 'v2.0'],
    ['view/extensions/emoji', 'v1.3.0'],
    ['view/extensions/likes', 'v1.3.0'],
    ['view/extensions/markdown', 'v1.3.0'],
    ['view/extensions/mentions', 'v1.3.1'],
    ['view/extensions/tags', 'v1.3.0'],
    ['view/flarum', 'v1.3.1'],
    ['view/locale/en', '9a03552'],
    ['view/locale/zh', '0ecda52'],
];

/*
git submodule foreach '
    if [[ "$path" =~ view/extensions/analytics ]]; then
        git fetch origin && git checkout 1.0.0
    fi

    if [[ "$path" =~ view/extensions/auth-github ]]; then
        git fetch origin && git checkout v0.1.0-beta.13
    fi
'
*/

function sync() {

    const commands = submodule_sync.map(function (item) {
        return `if [[ "$path" =~ ${item[0]} ]]; then
            git fetch origin && git checkout ${item[1]} -f
        fi
        `;
    });
    console.log(commands.join(''));

    execSync([
        'git', 'submodule', 'foreach', "'", commands.join(''), "'"
    ].join(' '), { stdio: 'inherit' })

    // submodule_sync.array.forEach(module => {
    //     ])
    // });
}

sync()

