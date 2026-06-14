# stew

`stew` is a terminal user interface (TUI) for the Salesforce CLI (`sf`), built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. It aims to provide a more interactive and visual way to manage Salesforce environments and metadata directly from the command line.

Stew was built with the goal of being an easy-to-integrate UI for interacting with the SF CLI in code editors other than VS Code, such as **Zed**, by leveraging its tasks feature.

You can checkout [SF Zed Tasks](https://github.com/brtheo/sf-zed-tasks) a comprehensive collection of `sf` and `stew` **Zed tasks**.
.

## Installation

Can be installed via npm (requires sudo on linux)

```bash
npm install -g @brtheo/stew
```

Clone the repository and build the binary:

```bash
git clone https://github.com/brtheo/stew.git
cd stew
./build.sh
```

## Features

### 1. Org Picker
Easily switch between your authenticated Salesforce orgs.
- Lists all authenticated orgs (Scratch Orgs, Sandboxes, Dev Hubs, etc.).
- Displays aliases, usernames, and instance URLs.
![Org Picker](./demo/org-picker.gif)

### 2. Metadata Generator
Streamline the creation of new Salesforce metadata with a simple, interactive form.
- **Lightning Web Components (LWC)**: Generate new LWC components.
- **Apex Classes**: Create new Apex classes.
- **Apex Triggers**: Create new Apex triggers with SObject selection.
- **Aura Components**: Generate new Aura components.
- **Aura Events**: Generate new Aura events.
- **Aura Applications**: Generate new Aura applications.
- Input validation to ensure metadata names are not empty.
- Real-time feedback on successful creation.
![Gen Metadata](./demo/gen-metadata.gif)

### 3. Metadata Retriever - Package XML generator
Browse, search, and select metadata from your Salesforce org to generate a package.xml file.
- Select a metadata type
- Browse all available metadata of that type and pick one or more items
- Generate a package.xml file with the selected metadata.
![Metadata Retriever](./demo/metadata-retriever.gif)

### 4. Authorize an Org
Wizard to authorize an org.

## Roadmap

- [*] gen aura related stuff
- [ ] gen vf page
- [*] gen project 
- [ ] gen test suite
- [*] authorize an org/dev hub
- [ ] get apex logs

## Prerequisites

- [Go](https://go.dev/doc/install) (1.19 or later)
- [Salesforce CLI (sf)](https://developer.salesforce.com/tools/sfdxcli) installed and authenticated with at least one org.
