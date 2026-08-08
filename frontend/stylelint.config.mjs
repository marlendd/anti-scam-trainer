import propertyGroups from 'stylelint-config-recess-order/groups'

export default {
    extends: ['stylelint-config-standard-scss'],

    plugins: ['stylelint-order'],

    rules: {
        'declaration-empty-line-before': null,

        'selector-class-pattern': '^(?:[a-z][a-zA-Z0-9]*|ag-[a-z0-9-]+)$',

        'selector-pseudo-class-no-unknown': [
            true,
            {
                ignorePseudoClasses: ['global'],
            },
        ],

        'order/properties-order': propertyGroups.map((group) => ({
            ...group,
            emptyLineBefore: 'always',
            noEmptyLineBetween: true,
        })),
    },
}
