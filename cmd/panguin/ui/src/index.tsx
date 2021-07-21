import { h, render } from '../deps/preact';
import { Explorer, ViewMode } from './components/explorer';

import "./index.css";

document.addEventListener('DOMContentLoaded', () => {
    const app = (
        <Explorer
            files={[
                { name: 'dummy file 1', size: 500, isDir: false, },
                { name: 'dummy file 2', size: 250, isDir: false, },
                { name: 'dummy file 3', size: 750, isDir: false, },
                { name: 'dummy file 4', size: 0, isDir: false, },
            ]}
            mode={ViewMode.Details} />
    );
    render(app, document.body);
});
