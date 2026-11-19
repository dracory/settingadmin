function initSettingManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.settingManagerUrls || {},
        settings: [],
        total: 0,
        page: 0,
        perPage: 50,
        loading: false,
        error: null,
        savingKey: null,
        filters: { key: '' }
      };
    },
    mounted() {
      this.loadSettings();
    },
    methods: {
      async loadSettings() {
        this.loading = true;
        this.error = null;
        try {
          const response = await fetch(this.urls.loadSettings, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              page: this.page,
              per_page: this.perPage,
              key: this.filters.key,
              conditions: [],
              sort_by: 'setting_key',
              sort: 'asc'
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            if (result.data.urls) {
              this.urls = { ...this.urls, ...result.data.urls };
            }
            this.settings = result.data.settings || [];
            this.total = result.data.total || 0;
          } else {
            this.error = result.message || 'Failed to load settings';
          }
        } catch (error) {
          console.error('Failed to load settings:', error);
          this.error = 'Failed to load settings';
        } finally {
          this.loading = false;
        }
      },
      settingLabel(key) {
        return String(key)
          .replace(/[._-]+/g, ' ')
          .replace(/\b\w/g, ch => ch.toUpperCase());
      },
      notifySuccess(message) {
        Notiflix.Notify.success(message, { position: 'right-top', timeout: 3000 });
      },
      notifyFailure(message) {
        Notiflix.Notify.failure(message, { position: 'right-top', timeout: 3000 });
      },
      async saveSetting(setting) {
        this.savingKey = setting.key;
        try {
          const response = await fetch(this.urls.saveSetting, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ key: setting.key, value: setting.value })
          });
          const result = await response.json();
          if (result.status === 'success') {
            this.notifySuccess('Setting saved: ' + setting.key);
          } else {
            this.notifyFailure(result.message || 'Failed to save setting');
          }
        } catch (error) {
          console.error('Failed to save setting:', error);
          this.notifyFailure('Failed to save setting');
        } finally {
          this.savingKey = null;
        }
      },
      async addSetting() {
        const { value: formValues } = await Swal.fire({
          title: 'Add Setting',
          html:
            '<input id="swal-key" class="swal2-input" placeholder="Key">' +
            '<input id="swal-value" class="swal2-input" placeholder="Value">',
          focusConfirm: false,
          showCancelButton: true,
          preConfirm: () => ({
            key: document.getElementById('swal-key').value.trim(),
            value: document.getElementById('swal-value').value
          })
        });
        if (!formValues || !formValues.key) return;
        await this.saveSetting({ key: formValues.key, value: formValues.value });
        await this.loadSettings();
      },
      async deleteSetting(setting) {
        const result = await Swal.fire({
          title: 'Delete setting?',
          text: setting.key,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: 'Delete',
          confirmButtonColor: '#dc3545'
        });
        if (!result.isConfirmed) return;
        try {
          const response = await fetch(this.urls.deleteSetting, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ setting_id: setting.id })
          });
          const json = await response.json();
          if (json.status === 'success') {
            this.notifySuccess('Setting deleted: ' + setting.key);
            await this.loadSettings();
          } else {
            this.notifyFailure(json.message || 'Failed to delete setting');
          }
        } catch (error) {
          console.error('Failed to delete setting:', error);
          this.notifyFailure('Failed to delete setting');
        }
      }
    }
  }).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initSettingManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initSettingManager();
      }
    });
  }
});
