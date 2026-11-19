function initSettingsEditor() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.settingAdminUrls || {},
        groups: [],
        loading: false,
        error: null,
        savingGroup: null,
        dirtyGroups: {},
        savedSnapshots: {}
      };
    },
    mounted() {
      this.loadSettings();
      window.addEventListener('beforeunload', this.beforeUnload);
    },
    unmounted() {
      window.removeEventListener('beforeunload', this.beforeUnload);
    },
    methods: {
      notifySuccess(message) {
        Notiflix.Notify.success(message, { position: 'right-top', timeout: 3000 });
      },
      notifyFailure(message) {
        Notiflix.Notify.failure(message, { position: 'right-top', timeout: 3000 });
      },

      // silent=true refetches without swapping to the loading spinner,
      // so the form stays mounted (used after a save to avoid scroll
      // jumps and flicker).
      async loadSettings(silent) {
        if (!silent) {
          this.loading = true;
          this.error = null;
        }
        try {
          const response = await fetch(this.urls.loadSettings, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({})
          });
          const result = await response.json();
          if (result.status === 'success') {
            this.groups = result.data.groups || [];
            this.groups.forEach(g => this.clearDirty(g));
          } else {
            if (!silent) {
              this.error = result.message || 'Failed to load settings';
            }
          }
        } catch (e) {
          console.error('Failed to load settings:', e);
          if (!silent) {
            this.error = 'Failed to load settings';
          }
        } finally {
          if (!silent) {
            this.loading = false;
          }
        }
      },

      // collectValues reads every named input inside the group card and
      // returns a {name: value} map. Checkboxes map to 'true'/'false';
      // radio buttons contribute only the checked value.
      collectValues(group) {
        const container = this.$refs['fields-' + group.id];
        const root = Array.isArray(container) ? container[0] : container;
        const values = {};
        if (!root) return values;

        root.querySelectorAll('[name]').forEach(function(el) {
          if (el.type === 'checkbox') {
            values[el.name] = el.checked ? 'true' : 'false';
          } else if (el.type === 'radio') {
            if (el.checked) values[el.name] = el.value;
          } else {
            values[el.name] = el.value;
          }
        });
        return values;
      },

      // snapshotValues returns a stable string of the group's current
      // values, used to detect unsaved edits.
      snapshotValues(group) {
        return JSON.stringify(this.collectValues(group));
      },

      // markDirty is called on input/change inside a group card. It
      // compares current inputs against the last-saved snapshot and
      // flags the group as dirty when they differ.
      markDirty(group) {
        // Defer to next tick so the DOM reflects the just-typed value.
        this.$nextTick(() => {
          const dirty = this.snapshotValues(group) !== this.savedSnapshots[group.id];
          this.dirtyGroups[group.id] = dirty;
        });
      },

      // clearDirty snapshots the group's persisted state (called after
      // load and after a successful save).
      clearDirty(group) {
        this.$nextTick(() => {
          this.savedSnapshots[group.id] = this.snapshotValues(group);
          this.dirtyGroups[group.id] = false;
        });
      },

      beforeUnload(e) {
        if (Object.values(this.dirtyGroups).some(Boolean)) {
          e.preventDefault();
        }
      },

      // refreshGroup refetches and splices in only the given group's
      // fields, leaving the rest of the DOM untouched.
      async refreshGroup(group) {
        try {
          const response = await fetch(this.urls.loadSettings, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({})
          });
          const result = await response.json();
          if (result.status === 'success') {
            const fresh = (result.data.groups || []).find(g => g.id === group.id);
            if (fresh) {
              group.fields = fresh.fields;
              this.clearDirty(group);
            }
          }
        } catch (e) {
          console.error('Failed to refresh group:', e);
        }
      },

      async saveGroup(group) {
        this.savingGroup = group.id;
        try {
          const response = await fetch(this.urls.saveGroup, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              group_id: group.id,
              values: this.collectValues(group)
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            this.notifySuccess('Settings saved: ' + group.title);
            // Refresh only this group's fields so other groups keep
            // their unsaved edits and scroll position is preserved.
            await this.refreshGroup(group);
          } else {
            this.notifyFailure(result.message || 'Failed to save settings');
          }
        } catch (e) {
          console.error('Failed to save settings:', e);
          this.notifyFailure('Failed to save settings');
        } finally {
          this.savingGroup = null;
        }
      }
    }
  }).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initSettingsEditor();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initSettingsEditor();
      }
    });
  }
});
