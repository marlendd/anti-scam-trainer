import { themeQuartz } from 'ag-grid-community'

export const appGridTheme = themeQuartz.withParams({
  backgroundColor: '#ffffff',
  foregroundColor: '#000000',

  fontFamily: 'inherit',
  fontSize: 14,
  headerFontWeight: 500,

  headerBackgroundColor: '#f5f5f7',
  oddRowBackgroundColor: '#fafafa',
  rowHoverColor: '#f0f1f3',

  borderColor: '#e5e7eb',
  wrapperBorder: true,
  wrapperBorderRadius: 20,

  rowHeight: 56,
  headerHeight: 48,
  paginationPanelHeight: 56,

  spacing: 10,

  browserColorScheme: 'inherit',

  chromeBackgroundColor: {
    ref: 'foregroundColor',
    mix: 0.07,
    onto: 'backgroundColor',
  },
})