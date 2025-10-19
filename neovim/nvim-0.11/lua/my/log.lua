return {
    _log_path = '/tmp/nvimdebug.log',
    _this_is_a_log_table = true,
    write = function(self, ...)
        if not self._this_is_a_log_table then
            error('Please call `log:write` with a semicolon')
        end
        if not self._fh then
            self._fh = io.open(self._log_path, 'a')
        end
        local buf = {}
        for i = 1, select('#', ...) do
            local v = select(i, ...)
            table.insert(buf, vim.inspect(v))
        end
        self._fh:write(
            os.date('%Y:%m:%d %H.%M.%S')
            .. ' '
            .. self.script_path()
            .. '\n'
            .. table.concat(buf, '\n')
            .. '\n\n'
        )
        self._fh:flush()
    end,
    script_path = function()
        return debug.getinfo(3, "S").source:sub(2)
    end,
}
