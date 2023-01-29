<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
    <xsl:template match="/">
    <html>
    <head>
      <meta charset="utf-8" />
      <title>home</title>
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <style>
      * {
        font-family: sans-serif;
        color: #888;
      }
      body {
        background-color: #111;
      }
      td, th {
        border-color: #888;
        border-style: solid;
      }
      th {
        color: #888;
        font-size: 150%;
        border-width: 0 0 2px 0;
      }
      td {
        color: #777;
        white-space: nowrap;
        border-width: 0;
      }
      a {
        color: #55f;
        text-decoration: none;
      }
      a:visited {
        color: #4d4;
      }
      a.dir {
        color: #97f;
      }
      a.dir:visited {
        color: #df4;
      }
      </style>
    </head>
    <body>
        <table border="0">
        <tr align="left">
            <th><a href="..">☚</a> home</th>
        </tr>
        <xsl:for-each select="list/*">
        <xsl:sort select="@name"/>

            <xsl:variable name="name">
                <xsl:value-of select="."/>
            </xsl:variable>
            <xsl:variable name="size">
                <xsl:if test="string-length(@size) &gt; 0">
                        <xsl:if test="number(@size) &gt; 0">
                            <xsl:choose>
                                    <xsl:when test="round(@size div 1024) &lt; 1"><xsl:value-of select="@size" /></xsl:when>
                                    <xsl:when test="round(@size div 1048576) &lt; 1"><xsl:value-of select="format-number((@size div 1024), '0.0')" />K</xsl:when>
                                    <xsl:otherwise><xsl:value-of select="format-number((@size div 1048576), '0.00')" />M</xsl:otherwise>
                            </xsl:choose>
                        </xsl:if>
                </xsl:if>
            </xsl:variable>
            <xsl:variable name="date">
                <xsl:value-of select="substring(@mtime,9,2)"/>-<xsl:value-of select="substring(@mtime,6,2)"/>-<xsl:value-of select="substring(@mtime,1,4)"/><xsl:text> </xsl:text>
                <xsl:value-of select="substring(@mtime,12,2)"/>:<xsl:value-of select="substring(@mtime,15,2)"/>:<xsl:value-of select="substring(@mtime,18,2)"/>
            </xsl:variable>

        <xsl:choose>
        <xsl:when test="string-length(@size) &gt; 0">
        <tr>
            <td><a href="{$name}"><xsl:value-of select="."/></a>
            <span style="float:right;"><xsl:value-of select="$size"/> | <xsl:value-of select="$date"/></span></td>
        </tr>
        </xsl:when>
        <xsl:otherwise>
        <tr>
            <td><a class="dir" href="{$name}"><xsl:value-of select="."/></a></td>
        </tr>
        </xsl:otherwise>
        </xsl:choose>

        </xsl:for-each>
        </table>
    </body>
    </html>
    </xsl:template>
</xsl:stylesheet>
